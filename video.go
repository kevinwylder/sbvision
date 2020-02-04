package sbvision

import (
	"net/http"
)

// Video is a generic video source
type Video struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	Thumbnail *Image  `json:"thumbnail"`
	Type      int64   `json:"type"`
	Duration  int64   `json:"duration"`
	FPS       float64 `json:"fps"`
	ClipCount int64   `json:"clips"`
}

// VideoList is a pagenated video lister interface
type VideoList interface {
	GetVideos(offset, count int) ([]Video, error)
	AddVideo(*Video, *Session) error
}

// VideoDownloadRequest is an incoming request to get the given video
type VideoDownloadRequest struct {
	Type int64  `json:"type"`
	URL  string `json:"url"`
}

// VideoHandler is able to "acquire" videos
type VideoHandler interface {
	HandleDownload(*VideoDownloadRequest) error
	HandleStream(http.ResponseWriter, *http.Request)
}

const (
	// YoutubeVideo is a video type that is streamed from youtube proxied through the server
	YoutubeVideo int64 = 1
)
