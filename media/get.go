package media

import "os"

// GetThumbnail opens the thumbnail file
func (sd *AssetDirectory) GetThumbnail(id int64) (*os.File, error) {
	return os.Open(sd.thumbnail(id))
}

// GetFrame opens the frame file
func (sd *AssetDirectory) GetFrame(id int64) (*os.File, error) {
	return os.Open(sd.frame(id))
}

// GetBound opens the bound file
func (sd *AssetDirectory) GetBound(id int64) (*os.File, error) {
	return os.Open(sd.bound(id))
}

// GetVideo opens the associted video file
func (sd *AssetDirectory) GetVideo(id int64) (*os.File, error) {
	return os.Open(sd.video(id))
}
