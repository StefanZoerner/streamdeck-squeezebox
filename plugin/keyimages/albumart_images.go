package keyimages

import (
	"errors"
	"fmt"
	"image"
	"net/http"
	"sync"
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
