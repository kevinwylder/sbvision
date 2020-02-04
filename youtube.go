package sbvision

import "time"

// YoutubeVideoInfo is the information required to serve a youtube video
type YoutubeVideoInfo struct {
	Video     *Video
	YoutubeID string
	MirrorURL string
	MirrorExp time.Time
}

// YoutubeVideoInfoStorage keeps track of indexed youtube links
type YoutubeVideoInfoStorage interface {
	GetYoutubeRecord(videoID int64) (*YoutubeVideoInfo, error)
	HasYoutubeRecord(youtubeID string) (*YoutubeVideoInfo, error)
	PutYoutubeRecord(*YoutubeVideoInfo) error
}
