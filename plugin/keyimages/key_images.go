package keyimages

import (
	"errors"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	"strconv"

	"github.com/golang/freetype/truetype"
	_ "golang.org/x/image/font"
	_ "golang.org/x/image/math/fixed"

	"github.com/nfnt/resize"
)

const (
	ALBUM_ART_1x1 = "1x1"
	ALBUM_ART_2x2 = "2x2"
	ALBUM_ART_3x3 = "3x3"
)

const (
	STREAMDECK_TILE_SIZE = 72
	STREAMDECK_GAP_SIZE  = 19
)

const (
	FontFile = "./assets/fonts/ArtsansC.ttf"
)

var (
	KeyForegroundColor = color.RGBA{0xd8, 0xd8, 0xd8, 255} // lightgray
	KeyBackgroundColor = color.RGBA{0x1d, 0x1d, 0x1f, 255} // darkgray
)

var (
	ttFont *truetype.Font
)

func init() {
	fontBytes, err := ioutil.ReadFile(FontFile)
	if err != nil {
		panic(err)
	}
	ttFont, err = truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
}

func ResizeAndCropImage(img image.Image, dimension string, tileNumber int) (image.Image, error) {
	switch dimension {
	case ALBUM_ART_1x1:
		return resizeAndCropImageFor1x1(img)
	case ALBUM_ART_2x2:
		return resizeAndCropImageFor2x2(img, tileNumber)
	case ALBUM_ART_3x3:
		return resizeAndCropImageFor3x3(img, tileNumber)
	default:
		err := errors.New("Unknown dimension for album art")
		return nil, err
	}
}

func resizeAndCropImageFor1x1(img image.Image) (image.Image, error) {

	var imgSize uint = STREAMDECK_TILE_SIZE
	smallerImg := resize.Resize(imgSize, imgSize, img, resize.Lanczos3)

	return smallerImg, nil
}

func resizeAndCropImageFor2x2(img image.Image, tileNumber int) (image.Image, error) {

	var imgSize uint = STREAMDECK_TILE_SIZE*2 + STREAMDECK_GAP_SIZE
	smallerImg := resize.Resize(imgSize, imgSize, img, resize.Lanczos3)

	var crop image.Image
	var err error

	switch tileNumber {
	case 1:
		crop, _ = CropImage(smallerImg, image.Rect(0, 0, STREAMDECK_TILE_SIZE-1, STREAMDECK_TILE_SIZE-1))
		break
	case 2:
		crop, _ = CropImage(smallerImg, image.Rect(72+19, 0, 163, 71))
		break
	case 3:
		crop, _ = CropImage(smallerImg, image.Rect(0, 72+19, 71, 163))
		break
	case 4:
		crop, _ = CropImage(smallerImg, image.Rect(72+19, 72+19, 163, 163))
		break
	default:
		crop = nil
		err = errors.New("Illegal tile number for 2x2: " + strconv.Itoa(tileNumber))
	}
	return crop, err
}

func resizeAndCropImageFor3x3(img image.Image, tileNumber int) (image.Image, error) {

	var imgSize uint = STREAMDECK_TILE_SIZE*3 + STREAMDECK_GAP_SIZE*2
	smallerImg := resize.Resize(imgSize, imgSize, img, resize.Lanczos3)

	var err error
	var crop image.Image

	switch tileNumber {
	case 1:
		crop, _ = CropImage(smallerImg, image.Rect(0, 0, 71, 71))
		break
	case 2:
		crop, _ = CropImage(smallerImg, image.Rect(72+19, 0, 72+19+72, 71))
		break
	case 3:
		crop, _ = CropImage(smallerImg, image.Rect(72+19+72+19, 0, 72+19+72+19+72, 71))
		break

	case 4:
		crop, _ = CropImage(smallerImg, image.Rect(0, 71+19, 71, 71+19+71))
		break
	case 5:
		crop, _ = CropImage(smallerImg, image.Rect(72+19, 71+19, 72+19+72, 71+19+71))
		break
	case 6:
		crop, _ = CropImage(smallerImg, image.Rect(72+19+72+19, 71+19, 72+19+72+19+72, 71+19+71))
		break

	case 7:
		crop, _ = CropImage(smallerImg, image.Rect(0, 71+19+71+19, 71, 72+19+72+19+72))
		break
	case 8:
		crop, _ = CropImage(smallerImg, image.Rect(72+19, 71+19+71+19, 72+19+72, 72+19+72+19+72))
		break
	case 9:
		crop, _ = CropImage(smallerImg, image.Rect(72+19+72+19, 71+19+71+19, 72+19+72+19+72, 72+19+72+19+72))
		break
	default:
		crop = nil
		err = errors.New("Illegal tile number for 3x3: " + strconv.Itoa(tileNumber))
	}

	return crop, err
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
