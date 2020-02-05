package sbvision

import (
	"io"
)

// Image is an object stored in s3
type Image string

// ImageUploader is an interface to upload images
type ImageUploader interface {
	UploadImage(io.Reader, string) (Image, error)
}

// ImageDownloader is an interface to download images
type ImageDownloader interface {
	DownloadImage(Image) (io.ReadCloser, error)
}

// ImageTracker adds images to the database
type ImageTracker interface {
	AddImage(Image, *Session) error
}
