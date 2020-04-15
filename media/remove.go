package media

import (
	"os"

	"github.com/kevinwylder/sbvision"
)

// RemoveVideo removes the given video by its id
func (sd *AssetDirectory) RemoveVideo(video *sbvision.Video) error {
	return os.Remove(sd.VideoPath(video))
}
