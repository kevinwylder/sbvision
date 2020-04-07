package sbvision

import (
	"io"
	"os"
)

// MediaStorage is an abstraction for the storage directory
type MediaStorage interface {
	GetBound(id int64) (*os.File, error)
	PutBound(id int64, data io.Reader) error

	GetFrame(id int64) (*os.File, error)
	PutFrame(id int64, data io.Reader) error

	GetThumbnail(id int64) (*os.File, error)
	PutThumbnail(id int64, data io.Reader) error
}
