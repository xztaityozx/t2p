package palette

import (
	"image"
	"image/color"
)

type Zoom struct {
	scale int
}

func NewZoom(s int) Zoom {
	return Zoom{scale: s}
}

// 1ドットをZoom.scale * Zoom.scaleのサイズにした画像を返す
// TODO: 計算量が最悪
func (z Zoom) ScaleUp(img *image.RGBA) *image.RGBA {
	var box [][]color.Color
	rct := img.Bounds()
	for y := rct.Min.Y; y < rct.Max.Y; y++ {
		var line []color.Color
		for x := rct.Min.X; x<rct.Max.X;x++ {
			for i := 0; i<z.scale;i++ {
				line = append(line, img.At(x, y))
			}
		}
		for i := 0; i<z.scale;i++ {
			box = append(box, line)
		}
	}

	newRct := image.Rect(0,0,len(box[0]),len(box))
	rt := image.NewRGBA(newRct)

	for y := newRct.Min.Y; y<newRct.Max.Y; y++ {
		for x := newRct.Min.X;x<newRct.Max.X; x++ {
			rt.Set(x,y,box[y][x])
		}
	}

	return rt
}
