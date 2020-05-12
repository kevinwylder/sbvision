package skateboard

import (
	"image"
	"image/color"
)

type Image []byte

func (data Image) At(x, y int) color.Color {
	i := 4 * (((499 - y) * 500) + x)
	if i < 0 || i >= len(data) {
		return color.Black
	}
	return color.RGBA{
		R: data[i],
		G: data[i+1],
		B: data[i+2],
		A: data[i+3],
	}
}

func (data Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, 500, 500)
}

func (data Image) ColorModel() color.Model {
	return color.RGBAModel
}
