package sbvision

import (
	"fmt"
	"io"
)

// KeyValueStore is an abstraction for files
type KeyValueStore interface {
	PutAsset(key Key, value io.Reader) error
	GetAsset(key Key) (io.ReadCloser, error)
}

// Key is a string that has been formatted by one of the functions below
type Key string

// Thumbnail gets the key for a video thumbnail to be used with KeyValueStore
func (video *Video) Thumbnail() Key {
	return Key(fmt.Sprintf("thumbnail/%d.jpg", video.ID))
}

// Key gets the key for an image frame to be used with KeyValueStore
func (frame *Frame) Key() Key {
	return Key(fmt.Sprintf("frame/%d.png", frame.ID))
}

// Key gets the key for a cropped image bound to be used with KeyValueStore
func (bounds *Bound) Key() Key {
	return Key(fmt.Sprintf("bound/%d.png", bounds.ID))
}
