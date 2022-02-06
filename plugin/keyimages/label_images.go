package keyimages

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"io/ioutil"
)

const (
	fontFile = "./assets/fonts/ArtsansC.ttf"
)

var (
	ttFont *truetype.Font
)

func init() {
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		panic(err)
	}
	ttFont, err = truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
}

func CreateKeyImageWithNumber(n int) image.Image {

	keyImage := createImage(72, 72, KeyBackgroundColor)
	fontSize := 48.

	d := &font.Drawer{
		Dst: &keyImage,
		Src: &image.Uniform{KeyForegroundColor},
		Face: truetype.NewFace(ttFont, &truetype.Options{
			Size:    fontSize,
			DPI:     72,
			Hinting: font.HintingNone,
		}),
	}

	if n > 99 {
		n = 99
	}

	if n < 0 {
		n = 0
	}

	if n < 10 {
		d.Dot = fixed.P(20, 50)
	} else {
		d.Dot = fixed.P(10, 50)
	}
	d.DrawString(fmt.Sprintf("%d", n))

	return &keyImage
}
