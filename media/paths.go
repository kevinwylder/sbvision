package media

import (
	"fmt"
	"path"
)

// AssetDirectory is a folder that holds images
type AssetDirectory struct {
	path string
}

func (sd *AssetDirectory) thumbnail(id int64) string {
	return path.Join(sd.path, fmt.Sprintf("thumbnail/%d.jpg", id))
}

func (sd *AssetDirectory) frame(id int64) string {
	return path.Join(sd.path, fmt.Sprintf("frame/%d.png", id))
}

func (sd *AssetDirectory) bound(id int64) string {
	return path.Join(sd.path, fmt.Sprintf("bound/%d.png", id))
}

func (sd *AssetDirectory) video(id int64) string {
	return path.Join(sd.path, fmt.Sprintf("video/%d.mp4", id))
}
