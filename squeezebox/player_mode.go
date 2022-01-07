package squeezebox

import (
	"fmt"
	"net"
	"strings"
)

func SetPlayerMode(player_id string, mode string) error {

	connection_string := fmt.Sprintf("%s:%d", hostname, port)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return err
	}
	defer con.Close()

	cmd := fmt.Sprintf("%s mode %s\n", player_id, mode);
	_, err = performCommand(con, cmd)

	return err
}

func GetPlayerMode(player_id string) (string, error) {

	connection_string := fmt.Sprintf("%s:%d", hostname, port)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return "", err
	}
	defer con.Close()

	cmd := fmt.Sprintf("%s mode ?\n", player_id);
	s, err := performCommand(con, cmd)
	if err != nil {
		return "", err
	}

	fmt.Printf("[%s]\n", s)

	tokens := strings.Split(s, " ")
	return  tokens[2], nil
}


func TogglePlayerMode(player_id string) (string, error) {

	connection_string := fmt.Sprintf("%s:%d", hostname, port)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return "", err
	}
	defer con.Close()

	cmd := fmt.Sprintf("%s mode ?", player_id);
	replyString, err := performCommand(con, cmd)
	if err != nil {
		return "", err
	}

	parts := strings.Split(replyString, " ")
	current_mode := parts[2]
	if strings.Count(current_mode, "play") == 1 {
		cmd := fmt.Sprintf("%s mode %s\n", player_id, "pause");
		_, err = performCommand(con, cmd)
		if err != nil {
			return "", err
		}
		return "pause", nil
	} else if strings.Count(current_mode, "pause") == 1 || strings.Count(current_mode, "stop") == 1 {
		cmd := fmt.Sprintf("%s mode %s\n", player_id, "play");
		_, err = performCommand(con, cmd)
		if err != nil {
			return "", err
		}
		return "play", nil
	}
	return "", nil
}