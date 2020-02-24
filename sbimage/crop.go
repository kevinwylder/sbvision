package sbimage

import (
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/kevinwylder/sbvision"
)

// PngCropper holds a decoded PNG and can get bounds of that image
type PngCropper struct {
	image image.Image
}

// Crop will fetch and decode the frame in this png
func Crop(data io.Reader) (*PngCropper, error) {
	image, err := png.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("\n\tError decoding png: %s", err.Error())
	}
	return &PngCropper{
		image: image,
	}, nil
}

// GetCroppedPng crops the frame and returns a PNG of just the given bounds
func (frame *PngCropper) GetCroppedPng(bounds *sbvision.Bound, dst io.Writer) error {
	rectangle := image.Rect(int(bounds.X), int(bounds.Y), int(bounds.X+bounds.Width), int(bounds.Y+bounds.Height))
	var subImage image.Image
	switch i := frame.image.(type) {
	case *image.RGBA:
		subImage = i.SubImage(rectangle)
	case *image.NRGBA:
		subImage = i.SubImage(rectangle)
	case *image.RGBA64:
		subImage = i.SubImage(rectangle)
	case *image.NRGBA64:
		subImage = i.SubImage(rectangle)
	default:
		return fmt.Errorf("\n\tError: image is of unknown type")
	}
	err := png.Encode(dst, subImage)
	if err != nil {
		return fmt.Errorf("\n\tError encoding subImage: %s", err.Error())
	}
	return nil
}
