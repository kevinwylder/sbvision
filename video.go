package sbvision

import (
	"net/http"
	"time"
)

// Video is a generic video source
type Video struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Format    string `json:"format"`
	Type      int64  `json:"type"`
	Duration  int64  `json:"duration"`
	ClipCount int64  `json:"clips"`
}

// VideoHandler is able to acquire videos
type VideoHandler interface {
	// discovery phase creates a video for the given URL and stores it in the database
	HandleDiscover(url string) (*Video, error)
	// the streaming phase gets part of the video for the HTML5 video tag
	HandleStream(http.ResponseWriter, *http.Request)
}

const (
	// YoutubeVideo is a video type that is streamed from youtube proxied through the server
	YoutubeVideo int64 = 1
)

// YoutubeVideoInfo is the database information on a youtube video
type YoutubeVideoInfo struct {
	VideoID   int64
	YoutubeID string
	MirrorURL string
	MirrorExp time.Time
}
