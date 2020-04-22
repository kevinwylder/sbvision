package media

import (
	"fmt"
	"path"
)

// Bound is a cropped frame around a skateboard
func (sd *AssetDirectory) Bound(id int64) string {
	return path.Join(sd.path, fmt.Sprintf("bound/%d.png", id))
}

// VideoPath is the path to the video given it's source URL
func (sd *AssetDirectory) VideoPath(id int64) string {
	return path.Join(sd.path, fmt.Sprintf("video/%d", id))
}
