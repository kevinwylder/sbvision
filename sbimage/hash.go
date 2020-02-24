package sbimage

import (
	"fmt"
	"hash/crc64"
	"image"
	"image/png"
	"io"
)

// HashImage computes the crc64 of the image data
func HashImage(data io.Reader) (int64, error) {

	img, err := png.Decode(data)
	if err != nil {
		return 0, fmt.Errorf("\n\tCould not decode PNG: %s", err.Error())
	}

	table := crc64.MakeTable(crc64.ECMA)

	var hash uint64
	switch i := img.(type) {
	case *image.RGBA:
		hash = crc64.Checksum(i.Pix, table)
	case *image.NRGBA:
		hash = crc64.Checksum(i.Pix, table)
	case *image.RGBA64:
		hash = crc64.Checksum(i.Pix, table)
	case *image.NRGBA64:
		hash = crc64.Checksum(i.Pix, table)
	default:
		return 0, fmt.Errorf("\n\tError: image is of unknown type")
	}

	return int64(hash), nil
}
