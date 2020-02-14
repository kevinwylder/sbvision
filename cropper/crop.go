package cropper

import (
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/kevinwylder/sbvision"
)

// PngCropper is able to crop png files
type PngCropper struct {
	assets sbvision.KeyValueStore
}

// NewPngCropper is the constructor for the png cropper
func NewPngCropper(assets sbvision.KeyValueStore) *PngCropper {
	return &PngCropper{
		assets: assets,
	}
}

// FrameCropper holds a decoded PNG and can get bounds of that image
type FrameCropper struct {
	image image.Image
	frame *sbvision.Frame
}

// GetFrame will fetch and decode the frame in this png
func (cropper *PngCropper) GetFrame(frame *sbvision.Frame) (*FrameCropper, error) {
	reader, err := cropper.assets.GetAsset(string(frame.Image))
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting asset: %s", err.Error())
	}
	image, err := png.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("\n\tError decoding png: %s", err.Error())
	}
	return &FrameCropper{
		image: image,
		frame: frame,
	}, nil
}

// GetCroppedPng crops the frame and returns a PNG of just the given bounds
func (frame *FrameCropper) GetCroppedPng(bounds *sbvision.Bound, dst io.Writer) error {
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
