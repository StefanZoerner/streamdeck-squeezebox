package keyimages

import (
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	"net/http"
	"strconv"
	"sync"
)

const (
	ALBUM_ART_1x1 = "1x1"
	ALBUM_ART_2x2 = "2x2"
	ALBUM_ART_3x3 = "3x3"
)

var (
	coverCache map[string]image.Image
	lock       sync.Mutex
)

func init() {
	coverCache = make(map[string]image.Image)
}

func GetImageByUrl(url string) (image.Image, error) {

	lock.Lock()
	defer lock.Unlock()

	if img, ok := coverCache[url]; ok {
		return img, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Status %d", res.StatusCode))
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	if err != nil {

		if img.Bounds().Max.X > 300 {
			rect := image.Rectangle{
				Min: image.Point{0, 300},
				Max: image.Point{0, 300},
			}
			img, err = CropImage(img, rect)
		}

		coverCache[url] = img
	}

	return img, err
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
		crop, _ = CropImage(smallerImg, image.Rect(0, 0, StreamdeckTileSize-1, StreamdeckTileSize-1))
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

	var imgSize uint = StreamdeckTileSize*3 + StreamdeckGapSize*2
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
