package main

import (
	"image"
	"image/color"

	"github.com/kevinwylder/sbvision"
)

func (rt *runtime) drawImageV1(in image.Image, box sbvision.Bound, maxDimension int64, rotation sbvision.Quaternion) image.Image {
	var (
		width  = 1000
		height = 500
		scale  = float64(maxDimension) * 1 / float64(height)
		cx     = float64(2*box.X+box.Width) / 2.0
		cy     = float64(2*box.Y+box.Height) / 2.0
	)

	skateboard := rt.skateboard.Render(rotation)
	rect := image.Rect(0, 0, width, height)
	out := image.NewRGBA(rect)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			x := cx + float64(i-250)*scale
			y := cy + float64(j-250)*scale
			k := 4 * (((499 - j) * 500) + i - 500)
			if i > 500 && skateboard[k+3] != 0 {
				out.Set(i, j, color.RGBA{
					R: skateboard[k+0],
					G: skateboard[k+1],
					B: skateboard[k+2],
					A: 255,
				})
			} else {
				out.Set(i, j, in.At(int(x+.5), int(y+.5)))
			}
		}
	}
	return out
}
