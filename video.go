package sbvision

// Video is a generic video source
type Video struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Thumbnail *Image `json:"thumbnail"`
	Type      string `json:"type"`
	Duration  int    `json:"duration"`
	FPS       int    `json:"fps"`
	ClipCount int    `json:"clips"`
}

// Frame is a frame of a video
type Frame struct {
	ID    int64
	Image *Image
	Video *Video `json:"video"`
	Time  int64  `json:"time"`
}

// VideoList is a pagenated video lister interface
type VideoList interface {
	GetVideos(offset, count int) ([]Video, error)
}
