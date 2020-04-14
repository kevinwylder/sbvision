package video

import (
	"os"

	"github.com/kevinwylder/sbvision/media/sources"

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

	request.setStatus("Starting Request")
	var err error
	request.source, err = sources.FindVideoSource(request.url)
	if err != nil {
		request.setStatus("Error finding video" + err.Error())
		return
	}

	request.setStatus("found " + request.source.Title())
	process, err := StartDownload(request.source.URL())
	if err != nil {
		request.setStatus(err.Error())
		return
	}

	// wait for the resolution, fps, and duration to be decoded
	progress := process.Progress()
	_, more := <-progress
	if !more {
		// bad sign, was there a problem with the source?
		if err := process.Error(); err != nil {
			request.setStatus("Error decoding source: " + err.Error())
			return
		}
	}

	request.setStatus("Adding to the database")
	process.Info.Title = request.source.Title()
	process.Info.Type = request.source.Type()
	err = q.database.AddVideo(&process.Info, request.user)
	if err != nil {
		request.setStatus("Failed to add data to the database - " + err.Error())
		process.Cancel()
		return
	}
	// remove video info from the database if not successful
	defer func() {
		if !request.WasSuccess {
			q.database.RemoveVideo(&process.Info)
		}
	}()

	request.setStatus("Getting thumbnail")
	data, err := request.source.GetThumbnail()
	if err != nil {
		request.setStatus("Error getting thumbnail - " + err.Error())
		return
	}
	err = q.assets.PutThumbnail(process.Info.ID, data)
	if err != nil {
		request.setStatus("Error storing thumbnail - " + err.Error())
		return
	}

	// wait for the rest of the video to be encoded, send pretty info to UI
	request.Info = &process.Info
	for {
		time, more := <-progress
		if time != "" {
			request.setStatus("Scanning Video - " + time + " of " + process.Info.Duration)
		}
		if !more {
			break
		}
	}
	if err = process.Error(); err != nil {
		request.setStatus(err.Error())
		return
	}

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

	request.setStatus("Complete")
	request.WasSuccess = true

}
