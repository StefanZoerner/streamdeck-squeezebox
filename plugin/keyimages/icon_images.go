package keyimages

import (
	"fmt"
	"github.com/samwho/streamdeck"
	"image"
	"os"
)

const (
	keyImageFilePath = "./assets/images/keys"
)

func GetStreamDeckImageForIcon(iconName string) (string, error)  {

	filename := fmt.Sprintf("%s/%s.png", keyImageFilePath, iconName)

	infile, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		return "", err
	}

	result, err := streamdeck.Image(img)
	if err != nil {
		return "", err
	}

	return result, nil
}