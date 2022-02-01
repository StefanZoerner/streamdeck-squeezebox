package squeezebox

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func CheckConnectionToPlayer(hostname string, port int, playerID string) error {

	connectionString := fmt.Sprintf("%s:%d", hostname, port)
	con, err := net.Dial("tcp", connectionString)
	if err != nil {
		return err
	}
	defer con.Close()

	cmd := fmt.Sprintf("%s connected ?\n", playerID)
	_, err = con.Write([]byte(cmd))
	if err != nil {
		return err
	}

	reply := make([]byte, 1024)
	_, err = con.Read(reply)
	if err != nil {
		return err
	}

	sReply := string(reply)
	sReply = strings.ReplaceAll(sReply, "\n", "")
	result := strings.Split(sReply, " ")[2]

	if strings.Contains(result, "%3F") {
		return errors.New("Player " + playerID + " not connected to server.")
	}

	return nil
}

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

	fmt.Println(response)

	artworkURL, _ := getTagValueFromResponseLine(response, "artworkURL")
	if artworkURL != "" {
		url = artworkURL
	} else {
		coverid, _ := getTagValueFromResponseLine(response, "coverid")
		if coverid != "" {
			// http://elfman:9002/music/1cec6e2c/cover.jpg
			url = fmt.Sprintf("http://%s:%d/music/%s/cover.jpg", cp.Hostname, cp.HTTPPort, coverid)
		}
	}

	return url, err
}
