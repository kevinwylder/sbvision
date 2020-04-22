package cdn

import (
	"fmt"
	"path"
)

// VideoDirectory returns the path for a video directory
func VideoDirectory(id int64) string {
	return fmt.Sprintf("/video/%d", id)
}

// VideoThumbnail returns the path for a video thumbnail
func VideoThumbnail(id int64) string {
	return path.Join(VideoDirectory(id), "thumbnail.jpg")
}

// VideoPlaylist returns the path for a video playlist
func VideoPlaylist(id int64) string {
	return path.Join(VideoDirectory(id), "playlist.m3u8")
}

// VideoMP4 returns the actual filepath for the whole video
func VideoMP4(id int64) string {
	return path.Join(VideoDirectory(id), "video.mp4")
}

// Bound is a path to the image represented by this bound
func Bound(id int64) string {
	return fmt.Sprintf("/bound/%d.png", id)
}
