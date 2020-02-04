package youtube

import (
	"fmt"
	"net/http"

	"github.com/kevinwylder/sbvision"
)

type youtubeHandler struct {
	db sbvision.YoutubeVideoInfoStorage
}

// NewYoutubeHandler constructs a namespace for downloading youtube video info
func NewYoutubeHandler(storage sbvision.YoutubeVideoInfoStorage) sbvision.VideoHandler {
	return &youtubeHandler{db: storage}
}

func (dl *youtubeHandler) HandleDownload(req *sbvision.VideoDownloadRequest) error {
	if req.Type != 1 {
		return fmt.Errorf("Download request is not a youtube type")
	}

}

func (dl *youtubeHandler) HandleStream(videoID int64, w http.ResponseWriter, r *http.Request) {

}
