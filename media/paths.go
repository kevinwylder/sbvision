package media

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"path"

	"github.com/kevinwylder/sbvision"
)

// AssetDirectory is a folder that holds stuff and should be put on a CDN
type AssetDirectory struct {
	path string
}

func (sd *AssetDirectory) thumbnail(id int64) string {
	return path.Join(sd.path, fmt.Sprintf("thumbnail/%d.jpg", id))
}

func (sd *AssetDirectory) bound(id int64) string {
	return path.Join(sd.path, fmt.Sprintf("bound/%d.png", id))
}

// VideoPath is the path to the video given it's source URL
func (sd *AssetDirectory) VideoPath(video *sbvision.Video) string {
	var data []byte
	switch video.Type {
	case sbvision.RedditVideo:
		data = []byte(video.SourceURL)
	case sbvision.YoutubeVideo:
		data = []byte(video.ShareURL)
	case sbvision.UploadedVideo:
		data = []byte(video.Title + video.UploadedAt)
	}
	sum := sha1.Sum(data)
	return path.Join(sd.path, fmt.Sprintf("video/%s.mp4", base64.URLEncoding.EncodeToString(sum[:])))
}
