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
	Width     int64     `json:"width"`
	Height    int64     `json:"height"`
	FPS       float64   `json:"fps"`
	Duration  string    `json:"duration"`
	Type      VideoType `json:"type"`
	ClipCount int64     `json:"clips"`
	URL       string    `json:"-"`
	OriginURL string    `json:"-"`
	LinkExp   time.Time `json:"-"`
}

// VideoType is an enum of subclasses of the video abstraction
type VideoType int64

const (
	// YoutubeVideo type means that this video has a record in the youtube_videos table
	YoutubeVideo VideoType = 1
	// RedditVideo type means that this video has a record in the reddit_videos table
	RedditVideo VideoType = 2
)

// VideoSource is an interface that wraps all video sources (youtube, reddit, file upload...)
type VideoSource interface {
	Title() string
	URL() string
	Type() VideoType
	GetThumbnail() (io.ReadCloser, error)
}
