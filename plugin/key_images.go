package plugin

import (
	"github.com/samwho/streamdeck"
	"image"
	"os"
)

func getImageByFilename(filename string) (string, error){
	infile, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer infile.Close()
	image, _, err := image.Decode(infile)

	result, err := streamdeck.Image(image)
	if err != nil {
		return "", err
	}

	return result, nil
}
