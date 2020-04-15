package media

import (
	"os"
	"sync"
	"time"

	"github.com/kevinwylder/sbvision"
)

// GetThumbnail opens the thumbnail file
func (sd *AssetDirectory) GetThumbnail(id int64) (*os.File, error) {
	return os.Open(sd.thumbnail(id))
}

// GetBound opens the bound file
func (sd *AssetDirectory) GetBound(id int64) (*os.File, error) {
	return os.Open(sd.bound(id))
}

// Ranger handles byte range reading. It is **supposed to be** thread safe
type Ranger struct {
	mutex sync.Mutex
	path  string
	file  *os.File
	close *time.Timer
}

// GetVideo gets a "ranger" to get byte ranges of the video.
func (sd *AssetDirectory) GetVideo(video *sbvision.Video) *Ranger {
	return &Ranger{
		path: sd.VideoPath(video),
	}
}

// GetRange gets a byte range of the file as specified by the "Range" http header.
func (r *Ranger) GetRange(start, end int64) ([]byte, error) {
	if end == 0 {
		end = start + 1024*1024
	}
	// ensure the file is open
	r.mutex.Lock()
	if r.file == nil {
		file, err := os.Open(r.path)
		if err != nil {
			r.mutex.Unlock()
			return nil, err
		}
		r.file = file
	}
	if r.close != nil {
		r.close.Stop()
	}
	r.close = time.AfterFunc(time.Minute, func() {
		if r.file == nil {
			return
		}
		r.mutex.Lock()
		r.file.Close()
		r.file = nil
		r.mutex.Unlock()
	})
	_, err := r.file.Seek(start, 0)
	if err != nil {
		r.mutex.Unlock()
		return nil, err
	}
	data := make([]byte, end-start)
	n, err := r.file.Read(data)
	r.mutex.Unlock()
	if err != nil {
		return nil, err
	}
	return data[0:n], nil

}
