package sbvision

import (
	"fmt"
	"io"
	"strings"
)

// Image is an object stored in s3
type Image struct {
	ID  int64
	Key string
}

// ImageManager is an interface to upload and download images
type ImageManager interface {
	UploadImage(io.Reader, string) (*Image, error)
	DownloadImage() (io.ReadCloser, error)
}

// ImageTracker adds IDs to the images and puts them in the database
type ImageTracker interface {
	TrackImage(*Image) error
}

// MarshalJSON overrides json marshalling so that the frontend just gets an ID
func (i *Image) MarshalJSON() ([]byte, error) {
	return []byte(`"/images/` + i.Key + `"`), nil
}

// UnmarshalJSON reverses the marshalling to get a struct out of the url
func (i *Image) UnmarshalJSON(data []byte) error {
	if len(data) < 11 || !strings.HasPrefix(string(data), `"/images/`) {
		return fmt.Errorf("cannot parse %s", string(data))
	}
	i.Key = string(data[10 : len(data)-1])
	return nil
}
