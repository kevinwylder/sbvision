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

// VideoPlaylist returns the .m3u8 path of the HLS content
func (sd *AssetDirectory) VideoPlaylist(id int64) string {
	return path.Join(sd.VideoPath(id), "playlist.m3u8")
}

// VideoFile is the path to the fallback file if HLS doesn't pan out
func (sd *AssetDirectory) VideoFile(id int64) string {
	return path.Join(sd.VideoPath(id), "video.mp4")
}
