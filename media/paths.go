package media

import (
	"fmt"
	"path"
)

// Thumbnail returns the path to the thumbnail directory
func (sd *AssetDirectory) Thumbnail(id int64) string {
	return path.Join(sd.path, fmt.Sprintf("thumbnail/%d.jpg", id))
}

// Bound is a cropped frame around a skateboard
func (sd *AssetDirectory) Bound(id int64) string {
	return path.Join(sd.path, fmt.Sprintf("bound/%d.png", id))
}

// VideoPath is the path to the video given it's source URL
func (sd *AssetDirectory) VideoPath(id int64) string {
	return path.Join(sd.path, fmt.Sprintf("video/%d", id))
}
