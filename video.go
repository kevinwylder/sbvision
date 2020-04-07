package sbvision

import (
	"io"
	"time"
)

// Video is a generic video source
type Video struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Format    string    `json:"format"`
	Type      VideoType `json:"type"`
	Duration  int64     `json:"duration"`
	ClipCount int64     `json:"clips"`
	URL       string    `json:"-"`
	OriginURL string    `json:"-"`
	LinkExp   time.Time `json:"-"`
}

// VideoType describes what the subclass of the video abstraction is
type VideoType int64

const (
	// YoutubeVideo type means that this video has a record in the youtube_videos table
	YoutubeVideo VideoType = 1
	// RedditVideo type means that this video has a record in the reddit_videos table
	RedditVideo VideoType = 2
)

// VideoSource is an interface that wraps both reddit and youtube
type VideoSource interface {
	Update(*Video)
	GetThumbnail() (io.ReadCloser, error)
}
