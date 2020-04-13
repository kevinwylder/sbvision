package media

import "os"

// RemoveVideo removes the given video by its id
func (sd *AssetDirectory) RemoveVideo(id int64) error {
	return os.Remove(sd.video(id))
}
