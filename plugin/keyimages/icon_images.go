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

var (
	keyIconCache map[string]string
)

func init() {
	keyIconCache = make(map[string]string)

	iconNames := []string{
		"album_art",
		"pause",
		"play",
		"play_pause",
		"track_prev",
		"track_next",
		"volume_up",
		"volume_down",
	}

	for _, name := range iconNames {
		sdImage, err := loadStreamDeckImageForIcon(name)
		if err != nil {
			fmt.Printf("icon %s not found, %s\n", name, err.Error())
		} else {
			keyIconCache[name] = sdImage
		}
	}
}

func GetStreamDeckImageForIcon(iconName string) (string, error) {
	var result string
	var err error

	result, ok := keyIconCache[iconName]
	if !ok {
		err = fmt.Errorf("icon %s is unknown", iconName)
	}

	return result, err
}

func loadStreamDeckImageForIcon(iconName string) (string, error) {

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
