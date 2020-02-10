package sbvision

import "time"

// YoutubeVideoInfo is the information required to serve a youtube video
type YoutubeVideoInfo struct {
	VideoID   int64
	YoutubeID string
	MirrorURL string
	MirrorExp time.Time
}

// YoutubeVideoTracker keeps track of indexed youtube links
type YoutubeVideoTracker interface {
	AddYoutubeRecord(*YoutubeVideoInfo) error
}

// YoutubeSearch is for finding videos that have already been discovered
type YoutubeSearch interface {
	GetYoutubeRecord(videoID int64) (*YoutubeVideoInfo, error)
}
