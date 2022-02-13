package keyimages

import (
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg" // for JPEG images
	"net/http"
	"strconv"
	"sync"
)

const (
	AlbumArt1x1 = "1x1"
	AlbumArt2x2 = "2x2"
	AlbumArt3x3 = "3x3"

	cacheImageSize = 300
)

var (
	coverCache   map[string]image.Image
	lock         sync.Mutex
	defaultImage image.Image
)

func init() {
	var err error
	coverCache = make(map[string]image.Image)
	defaultImage, err = GetImageByFilename("./assets/images/album_art_default.png")
	if err != nil {
		panic(err)
	}
}

func GetImageByUrl(url string) (image.Image, error) {

	var img image.Image
	var err error

	lock.Lock()
	defer lock.Unlock()

	if img, ok := coverCache[url]; ok {
		return img, nil
	}

	if url == "" {
		img = defaultImage
	} else {
		res, err := http.Get(url)
		if err != nil {
			return defaultImage, err
		}
		if res.StatusCode != 200 {
			return defaultImage, fmt.Errorf("HTTP Status %d", res.StatusCode)
		}
		defer res.Body.Close()

		img, _, err = image.Decode(res.Body)
		if err != nil {
			return defaultImage, err
		}
	}

	if img != nil {
		if img.Bounds().Max.X > cacheImageSize {
			img = resize.Resize(cacheImageSize, cacheImageSize, img, resize.Lanczos3)
		}
		coverCache[url] = img
	}

	return img, err
}

func ResizeAndCropImage(img image.Image, dimension string, tileNumber int) (image.Image, error) {
	switch dimension {
	case AlbumArt1x1:
		return resizeAndCropImageFor1x1(img)
	case AlbumArt2x2:
		return resizeAndCropImageFor2x2(img, tileNumber)
	case AlbumArt3x3:
		return resizeAndCropImageFor3x3(img, tileNumber)
	default:
		err := errors.New("unknown dimension for album art")
		return nil, err
	}
}

func resizeAndCropImageFor1x1(img image.Image) (image.Image, error) {

	var imgSize uint = StreamdeckTileSize
	smallerImg := resize.Resize(imgSize, imgSize, img, resize.Lanczos3)

	return smallerImg, nil
}

func resizeAndCropImageFor2x2(img image.Image, tileNumber int) (image.Image, error) {

	var imgSize uint = StreamdeckTileSize*2 + StreamdeckGapSize
	smallerImg := resize.Resize(imgSize, imgSize, img, resize.Lanczos3)

	var crop image.Image
	var err error

	switch tileNumber {
	case 1:
		crop, _ = cropImage(smallerImg, image.Rect(0, 0, StreamdeckTileSize-1, StreamdeckTileSize-1))
	case 2:
		crop, _ = cropImage(smallerImg, image.Rect(72+19, 0, 163, 71))
	case 3:
		crop, _ = cropImage(smallerImg, image.Rect(0, 72+19, 71, 163))
	case 4:
		crop, _ = cropImage(smallerImg, image.Rect(72+19, 72+19, 163, 163))
	default:
		crop = nil
		err = errors.New("Illegal tile number for 2x2: " + strconv.Itoa(tileNumber))
	}
	return crop, err
}

func resizeAndCropImageFor3x3(img image.Image, tileNumber int) (image.Image, error) {

	var imgSize uint = StreamdeckTileSize*3 + StreamdeckGapSize*2
	smallerImg := resize.Resize(imgSize, imgSize, img, resize.Lanczos3)

	var err error
	var crop image.Image

	switch tileNumber {
	case 1:
		crop, _ = cropImage(smallerImg, image.Rect(0, 0, 71, 71))
	case 2:
		crop, _ = cropImage(smallerImg, image.Rect(72+19, 0, 72+19+72, 71))
	case 3:
		crop, _ = cropImage(smallerImg, image.Rect(72+19+72+19, 0, 72+19+72+19+72, 71))
	case 4:
		crop, _ = cropImage(smallerImg, image.Rect(0, 71+19, 71, 71+19+71))
	case 5:
		crop, _ = cropImage(smallerImg, image.Rect(72+19, 71+19, 72+19+72, 71+19+71))
	case 6:
		crop, _ = cropImage(smallerImg, image.Rect(72+19+72+19, 71+19, 72+19+72+19+72, 71+19+71))
	case 7:
		crop, _ = cropImage(smallerImg, image.Rect(0, 71+19+71+19, 71, 72+19+72+19+72))
	case 8:
		crop, _ = cropImage(smallerImg, image.Rect(72+19, 71+19+71+19, 72+19+72, 72+19+72+19+72))
	case 9:
		crop, _ = cropImage(smallerImg, image.Rect(72+19+72+19, 71+19+71+19, 72+19+72+19+72, 72+19+72+19+72))
	default:
		crop = nil
		err = errors.New("Illegal tile number for 3x3: " + strconv.Itoa(tileNumber))
	}

	return crop, err
}
