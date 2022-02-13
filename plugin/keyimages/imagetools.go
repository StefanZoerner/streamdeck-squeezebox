package keyimages

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"os"
)

const (
	StreamdeckTileSize = 72
	StreamdeckGapSize  = 19
)

var (
	KeyForegroundColor = color.RGBA{0xd8, 0xd8, 0xd8, 255} // lightgray
	KeyBackgroundColor = color.RGBA{0x1d, 0x1d, 0x1f, 255} // darkgray
)

func GetImageByFilename(filename string) (image.Image, error) {

	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func createImage(width, height int, c color.Color) image.RGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	fillRectangle(img, 0, 0, width, height, c)

	return *img
}

func fillRectangle(img *image.RGBA, startX, startY, width, height int, c color.Color) {
	rect := image.Rectangle{
		Min: image.Point{startX, startY},
		Max: image.Point{startX + width, startY + height},
	}
	draw.Draw(img, rect, &image.Uniform{c}, image.Point{}, draw.Src)
}

// Inspired by https://riptutorial.com/go/example/31687/cropping-image

func cropImage(img image.Image, cropRect image.Rectangle) (cropImg image.Image, err error) {

	//Interface for asserting whether `img`
	//implements SubImage or not.
	//This can be defined globally.
	type CropableImage interface {
		image.Image
		SubImage(r image.Rectangle) image.Image
	}

	if p, ok := img.(CropableImage); ok {
		// Call SubImage. This should be fast,
		// since SubImage (usually) shares underlying pixel.
		cropImg = p.SubImage(cropRect)
	} else if cropRect = cropRect.Intersect(img.Bounds()); !cropRect.Empty() {
		// If `img` does not implement `SubImage`,
		// copy (and silently convert) the image portion to RGBA image.
		rgbaImg := image.NewRGBA(cropRect)
		for y := cropRect.Min.Y; y < cropRect.Max.Y; y++ {
			for x := cropRect.Min.X; x < cropRect.Max.X; x++ {
				rgbaImg.Set(x, y, img.At(x, y))
			}
		}
		cropImg = rgbaImg
	} else {
		// Return an empty RGBA image
		cropImg = &image.RGBA{}
		err = errors.New("Cropping failed")
	}

	return cropImg, err
}
