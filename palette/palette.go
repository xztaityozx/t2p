package palette

import (
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// base64エンコードされたPNGのパレット画像
// 横幅16px、縦幅8pxである
// エンコード前の元画像はt2p/_image/table.png
const srcImg = `iVBORw0KGgoAAAANSUhEUgAAABAAAAAICAYAAADwdn+XAAAAAXNSR0IArs4c6QAAAe9JREFUKJEF
wUFIUwEAgOF/brahbbjmSmaa7pm1uWmZhFlJOBADkUVBnqrDqA4WhIFX6VJQBwkraBaSYFoY0UUP
09EUNK18a63U5Vwb6tRNXLqKTXt9nwyQjEVVCCoVxdokpoIs7PlrpCrOoBeqUFd08XzTgHc2xKdv
teT636MMJwkl9KjlJ5EduXJHKtNcRmUtZLtykRZNHlUlcrreDLHiFVmMRdiIRSmIVXN8ow/rwaP8
K7nI0P4O0jNLKFZdr7BmB2gSf1MbVpBjL6QvkeZZ1hbxTT2Kj5nsOdRG1GFl5+wNGjJ36HZ9Ztx1
innRi6y1vlGqzjFhy1WiPW2CE7OgbKY1uYhHlLhuqcMz8pPAcISQx82ttQkqdXEOCxnkO2qQn2+5
1z7s+4JzcoJtaY7slIEsixZ15C6rY4M4Vl+w7XzLNWGLdFqJrP8STwwC95fCdHa7kUmSJA386eSp
OMb0qBZ1b5CicgHhmBGzqZjy0iLWf+0iNBdnyj1JWJxjfNDMut9HwvcDuW2hp935cItgf5QD6QQX
dHn4VQY0GT6cbT00KoOUvhyh17KbMuEdHftSPNY1YZMF0NbfRGFssHP1QxCXKDH1dQMDmwzoouyt
XEbxvRbHwhrTyb8sP3pNylzHaPM5bldYsa7UMPNgnv9Ye8bsLd6zuwAAAABJRU5ErkJggg==`

type Palette struct {
	table []color.Color
}

func NewPalette(path string) Palette {
	return Palette{
		table: getColorTable(path),
	}
}

// getColorTable は色パレット画像からカラーテーブルを作って返す。
// 画像ファイルのパスに空文字列を指定すると組み込みの画像データから
// カラーテーブルを作って返す。
//
// なお、このカラーテーブルを作るのに使用可能な画像フォーマットはPNGのみである。
func getColorTable(path string) []color.Color {
	var r io.Reader
	if len(path) == 0 {
		r = base64.NewDecoder(base64.StdEncoding, strings.NewReader(srcImg))
	} else {
		var err error
		r, err = os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			logrus.WithError(err).Fatal("Failed open file: ", path)
		}
	}

	d, err := png.Decode(r)
	if err != nil {
		logrus.Fatal("Failed decode table image")
	}

	bounds := d.Bounds()

	var rt []color.Color

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rt = append(rt, d.At(x, y))
		}
	}

	// スペースは必ず(0,0,0,0)
	rt[int(' ')] = color.RGBA{R: 0, G: 0, B: 0, A: 0}

	return rt
}

func (p Palette) stringToRGBA(s string) []color.Color {
	var rt []color.Color

	for _, v := range s {
		rt = append(rt, p.table[int(v)])
	}

	return rt
}

// 1行の文字列を幅wドット、高さ1ドットの画像で返す
func (p Palette) toImage(s string, w int) []color.Color {
	var rt []color.Color
	img := p.stringToRGBA(s)

	for x := 0; x < w; x++ {
		if x >= len(s) {
			rt = append(rt, color.RGBA{R: 0, G: 0, B: 0, A: 0})
		} else if rune(x) > 127 {
			rt = append(rt, color.RGBA{R: 0, G: 0, B: 0, A: 255})
		} else {
			rt = append(rt, img[x])
		}
	}

	return rt
}

func (p Palette) Create(w, h int, src []string) *image.RGBA {
	rt := image.NewRGBA(image.Rect(0, 0, w, h))

	for y, v := range src {
		l := p.toImage(v, w)
		for x := 0; x < w; x++ {
			rt.Set(x, y, l[x])
		}
	}

	return rt
}

func (p Palette) GetPaletteSize() int {
	return len(p.table)
}
