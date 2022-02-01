package squeezebox

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func CheckConnectionToPlayer(hostname string, port int, player_id string) error {

	connection_string := fmt.Sprintf("%s:%d", hostname, port)

	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return err
	}

	defer con.Close()
	cmd := fmt.Sprintf("%s connected ?\n", player_id)
	_, err = con.Write([]byte(cmd))
	if err != nil {
		return err
	}

	reply := make([]byte, 1024)
	_, err = con.Read(reply)
	if err != nil {
		return err
	}

	s_reply := string(reply)
	s_reply = strings.ReplaceAll(s_reply, "\n", "")
	result := strings.Split(s_reply, " ")[2]

	if strings.Contains(result, "%3F") {
		return errors.New("Player " + player_id + " not connected to server.")
	}

	return nil
}

func GetCurrentArtworkUrl(cp ConnectionProperties, player_id string) (string, error) {

	url := ""

	connection_string := fmt.Sprintf("%s:%d", cp.Hostname, cp.CLIPort)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return "", err
	}
	defer con.Close()

	cmd := fmt.Sprintf("%s status - 1 tags:K,c\n", player_id)
	response, err := performCommand(con, cmd)
	if err != nil {
		return "", err
	}

	fmt.Println(response)

	artwork_url, _ := getTagValueFromResponseLine(response, "artwork_url")
	if artwork_url != "" {
		url = artwork_url
	} else {
		coverid, _ := getTagValueFromResponseLine(response, "coverid")
		if coverid != "" {
			// http://elfman:9002/music/1cec6e2c/cover.jpg
			url = fmt.Sprintf("http://%s:%d/music/%s/cover.jpg", cp.Hostname, cp.HTTPPort, coverid)
		}
	}

	return url, err
}
