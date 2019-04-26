package palette

import (
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"image"
	"image/color"
	"strings"
)

type Palette struct {
	table []color.RGBA
}

func NewPalette() Palette {
	return Palette{
		table: getColorTable(),
	}
}

func getColorTable() []color.RGBA {
	r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(srcImg))
	img, _, err := image.Decode(r)
	if err != nil {
		logrus.Fatal("Failed decode table image")
	}
	bounds := img.Bounds()

	var rt []color.RGBA
	box := image.NewRGBA(bounds)

	for y:=bounds.Min.Y;y<bounds.Max.Y;y++ {
		for x:= bounds.Min.X;x<bounds.Max.X;x++{
			rt = append(rt, box.RGBAAt(x,y))
		}
	}

	logrus.Info(bounds.Max)
	logrus.Info(len(rt))

	return rt
}

func (p Palette) stringToRGBA(s string) []color.RGBA {
	var rt []color.RGBA

	for _, v := range s {
		rt = append(rt, p.table[int(v)])
	}

	return rt
}

func (p Palette) ToImage(s string, w, h int) *image.RGBA {
	rt := image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}})
	img := p.stringToRGBA(s)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w+x
			if idx >= len(s) {
				rt.SetRGBA(x,y,color.RGBA{R:0,G:0,B:0,A:255})
			} else {
				rt.SetRGBA(x, y, img[y*w+x])
			}
		}
	}

	return rt
}

const srcImg = `
iVBORw0KGgoAAAANSUhEUgAAABAAAAAICAYAAADwdn+XAAAAAXNSR0IArs4c6QAAAe9JREFUKJEF
wUFIUwEAgOF/brahbbjmSmaa7pm1uWmZhFlJOBADkUVBnqrDqA4WhIFX6VJQBwkraBaSYFoY0UUP
09EUNK18a63U5Vwb6tRNXLqKTXt9nwyQjEVVCCoVxdokpoIs7PlrpCrOoBeqUFd08XzTgHc2xKdv
teT636MMJwkl9KjlJ5EduXJHKtNcRmUtZLtykRZNHlUlcrreDLHiFVmMRdiIRSmIVXN8ow/rwaP8
K7nI0P4O0jNLKFZdr7BmB2gSf1MbVpBjL6QvkeZZ1hbxTT2Kj5nsOdRG1GFl5+wNGjJ36HZ9Ztx1
innRi6y1vlGqzjFhy1WiPW2CE7OgbKY1uYhHlLhuqcMz8pPAcISQx82ttQkqdXEOCxnkO2qQn2+5
1z7s+4JzcoJtaY7slIEsixZ15C6rY4M4Vl+w7XzLNWGLdFqJrP8STwwC95fCdHa7kUmSJA386eSp
OMb0qBZ1b5CicgHhmBGzqZjy0iLWf+0iNBdnyj1JWJxjfNDMut9HwvcDuW2hp935cItgf5QD6QQX
dHn4VQY0GT6cbT00KoOUvhyh17KbMuEdHftSPNY1YZMF0NbfRGFssHP1QxCXKDH1dQMDmwzoouyt
XEbxvRbHwhrTyb8sP3pNylzHaPM5bldYsa7UMPNgnv9Ye8bsLd6zuwAAAABJRU5ErkJggg==
`
