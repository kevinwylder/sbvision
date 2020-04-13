package video

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/kevinwylder/sbvision"
)

// ProcessRequest is a request to process the given video source
type ProcessRequest struct {
	User       *sbvision.User
	Source     sbvision.VideoSource
	ID         string `json:"id"`
	Status     string `json:"status"`
	IsComplete bool   `json:"complete"`
	WasSuccess bool   `json:"success"`
}

// Database is the interface to persistently store a video
type Database interface {
	AddVideo(video *sbvision.Video) error
	RemoveVideo(video *sbvision.Video) error
}

// ProcessQueue is a queue of pending or processing requests
type ProcessQueue struct {
	assets   sbvision.MediaStorage
	database Database

	requests map[string]*ProcessRequest
	queue    chan *ProcessRequest
}

// Enqueue adds this source to the list of videos to process
func (q *ProcessQueue) Enqueue(user *sbvision.User, source sbvision.VideoSource) (*ProcessRequest, error) {
	data := make([]byte, 10)
	rand.Reader.Read(data)
	request := &ProcessRequest{
		User:   user,
		Source: source,
		ID:     base64.URLEncoding.EncodeToString(data),
	}

	select {
	case q.queue <- request:
		q.requests[request.ID] = request
		return request, nil
	default:
		return nil, fmt.Errorf("Queue is full, please come back later")
	}
}

// NewProcessQueue creates a queue wrapping the single ffmpeg process allowed
func NewProcessQueue(assets sbvision.MediaStorage, database Database) *ProcessQueue {
	queue := ProcessQueue{
		assets:   assets,
		database: database,
		queue:    make(chan *ProcessRequest, 20),
		requests: make(map[string]*ProcessRequest),
	}
	go queue.start()
	return &queue
}

func (q *ProcessQueue) start() {
	for {
		request := <-q.queue
		q.processRequest(request)
	}
}

func (q *ProcessQueue) processRequest(request *ProcessRequest) {
	defer func() {
		request.finish(q)
	}()

	process, err := StartDownload(request.Source.URL())
	if err != nil {
		request.setStatus(err.Error())
		return
	}

	progress := process.Progress()
	for {
		text, finished := <-progress
		if finished {
			break
		} else {
			request.setStatus("Downloading Video - " + text)
		}
	}

	err = process.Error()
	if err != nil {
		request.setStatus(err.Error())
		return
	}

	request.setStatus("Adding to the database")
	process.Info.Title = request.Source.Title()
	process.Info.Type = request.Source.Type()
	process.Info.OriginURL = request.Source.URL()
	err = q.database.AddVideo(&process.Info)
	if err != nil {
		request.setStatus("Failed to add data to the database - " + err.Error())
		return
	}
	// remove video info if not successful
	defer func() {
		if !request.WasSuccess {
			q.database.RemoveVideo(&process.Info)
		}
	}()

	request.setStatus("Storing Video")
	file, err := os.Open(process.OutputPath)
	if err != nil {
		request.setStatus("Failed to open video file")
		return
	}
	err = q.assets.PutVideo(process.Info.ID, file)
	if err != nil {
		request.setStatus("Failed to store video file")
		return
	}
	file.Close()
	// remove video file if not successful
	defer func() {
		if !request.WasSuccess {
			q.assets.RemoveVideo(process.Info.ID)
		}
	}()

	request.setStatus("getting thumbnail")
	data, err := request.Source.GetThumbnail()
	if err != nil {
		request.setStatus("Error getting thumbnail - " + err.Error())
		return
	}
	err = q.assets.PutThumbnail(process.Info.ID, data)
	if err != nil {
		request.setStatus("Error storing thumbnail - " + err.Error())
		return
	}

	request.setStatus("Complete")
	request.WasSuccess = true

}

func (r *ProcessRequest) finish(q *ProcessQueue) {
	r.IsComplete = true
	defer func() {
		time.Sleep(time.Minute)
		delete(q.requests, r.ID)
	}()
}

func (r *ProcessRequest) setStatus(status string) {
	r.Status = status
}

// Find looks up the request for the given id
func (q *ProcessQueue) Find(id string) (*ProcessRequest, bool) {
	request, exists := q.requests[id]
	return request, exists
}
