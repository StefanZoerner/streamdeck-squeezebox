package squeezebox

import (
	"fmt"
	"net"
)

func GetCurrentArtworkURL(cp ConnectionProperties, playerID string) (string, error) {

	url := ""

	connectionString := fmt.Sprintf("%s:%d", cp.Hostname, cp.CLIPort)
	con, err := net.Dial("tcp", connectionString)
	if err != nil {
		return "", err
	}
	defer con.Close()

	cmd := fmt.Sprintf("%s status - 1 tags:K,c\n", playerID)
	response, err := performCommand(con, cmd)
	if err != nil {
		return "", err
	}

	artworkURL, _ := getTagValueFromResponseLine(response, "artwork_url")
	if artworkURL != "" {
		url = artworkURL
	} else {
		coverid, _ := getTagValueFromResponseLine(response, "coverid")
		if coverid != "" {
			url = fmt.Sprintf("http://%s:%d/music/%s/cover.jpg", cp.Hostname, cp.HTTPPort, coverid)
		}
	}

	return url, err
}
