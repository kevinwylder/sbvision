package sbvision

import "io"

// User comes from the cognito user pool
type User struct {
	ID       int64
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Video is a generic video source
type Video struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Width      int64     `json:"width"`
	Height     int64     `json:"height"`
	FPS        float64   `json:"fps"`
	Duration   string    `json:"duration"`
	Type       VideoType `json:"type"`
	UploadedAt string    `json:"uploaded_at"`
	// SourceURL is the url that the video was downloaded from, or a path if the file was uploaded
	SourceURL string `json:"-"`
	// ShareURL is an optional url that describes where to find the video in a web browser
	ShareURL string `json:"-"`
}

// Clip is part of a video that has a trick
type Clip struct {
	ID       int64    `json:"id"`
	VideoID  int64    `json:"video"`
	Start    int64    `json:"start"`
	End      int64    `json:"end"`
	Username string   `json:"clipped_by"`
	Tricks   []string `json:"tricks"`
	Bounds   []Bound  `json:"bounds"`
}

// Bound is an area on a frame
type Bound struct {
	ID        int64      `json:"id"`
	Frame     int64      `json:"frame"`
	Width     int64      `json:"width"`
	Height    int64      `json:"height"`
	X         int64      `json:"x"`
	Y         int64      `json:"y"`
	Rotations []Rotation `json:"rotations"`
}

// Rotation is the angle that a user has voted for a bound
type Rotation struct {
	ID      int64   `json:"-"`
	BoundID int64   `json:"-"`
	R       float64 `json:"r"`
	I       float64 `json:"i"`
	J       float64 `json:"j"`
	K       float64 `json:"k"`
}

// VideoType is an enum of subclasses of the video abstraction
type VideoType int64

const (
	// YoutubeVideo type means that this video has a record in the youtube_videos table
	YoutubeVideo VideoType = 1
	// RedditVideo type means that this video has a record in the reddit_videos table
	RedditVideo VideoType = 2
	// UploadedVideo type means this video was uploaded
	UploadedVideo VideoType = 3
)

// VideoSource is an interface that wraps all video sources (youtube, reddit, file upload...)
type VideoSource interface {
	GetVideo() Video
	GetThumbnail() (io.ReadCloser, error)
}
