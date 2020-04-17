package video

import (
	"io"
	"os"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/media"
)

// Database is the interface to persistently store a video
type Database interface {
	AddVideo(video *sbvision.Video, user *sbvision.User) error
	RemoveVideo(video *sbvision.Video) error
}

// ProcessQueue is a queue of pending or processing requests
type ProcessQueue struct {
	assets   *media.AssetDirectory
	database Database

	requests map[int64]*ProcessRequest
	queue    chan *ProcessRequest
}

// NewProcessQueue creates a queue wrapping the single ffmpeg process allowed
func NewProcessQueue(assets *media.AssetDirectory, database Database) *ProcessQueue {
	queue := ProcessQueue{
		assets:   assets,
		database: database,
		queue:    make(chan *ProcessRequest, 20),
		requests: make(map[int64]*ProcessRequest),
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

	defer func() {
		if !request.WasSuccess {
			q.database.RemoveVideo(request.Info)
		}
	}()

	request.setStatus("Getting source information")
	if err := request.getSourceInformation(); err != nil {
		request.setStatus("Error getting source information: " + err.Error())
		return
	}

	request.setStatus("Getting video information")
	if err := request.getVideoInformation(); err != nil {
		request.setStatus("Error getting video info: " + err.Error())
		return
	}

	request.setStatus("Adding to the database")
	if err := request.addToDatabase(); err != nil {
		request.setStatus("Failed to add to the database - " + err.Error())
		return
	}

	request.setStatus("Getting Thumbnail")
	if err := request.getThumbnail(); err != nil {
		request.setStatus("Failed to get thumbnail - " + err.Error())
		return
	}

	request.setStatus("Processing Video")
	if err := request.processVideo(); err != nil {
		request.setStatus("Error processing video - " + err.Error())
	}

	request.setStatus("Complete")
	request.WasSuccess = true

}

func (r *ProcessRequest) getSourceInformation() error {
	source, err := r.getSource()
	if err != nil {
		return err
	}
	r.source = source
	info := source.GetVideo()
	r.Info = &info
	return nil
}

func (r *ProcessRequest) getVideoInformation() error {
	process := getInfo(r.Info)
	for range <-process.Progress() {
	}
	return process.Error()
}

func (r *ProcessRequest) addToDatabase() error {
	return r.q.database.AddVideo(r.Info, r.user)
}

func (r *ProcessRequest) getThumbnail() error {
	data, err := r.source.GetThumbnail()
	if err != nil {
		return err
	}
	file, err := os.Create(r.q.assets.Thumbnail(r.Info.ID))
	if err != nil {
		return err
	}
	_, err = io.Copy(file, data)
	return err
}

func (r *ProcessRequest) processVideo() error {
	process := r.q.startDownload(r.Info)
	for status := range process.Progress() {
		r.setStatus(status)
	}
	return process.Error()
}
