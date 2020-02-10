package sbvision

import "io"

// KeyValueStore is an abstraction for files
type KeyValueStore interface {
	PutAsset(key string, value io.Reader) error
	GetAsset(key string) (io.ReadCloser, error)
}
