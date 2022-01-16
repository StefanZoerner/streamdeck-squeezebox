package squeezebox

import (
	"fmt"
	"net"
)

// SetPlayerMode sets the mode of the player with player_id to mode.
// Possible values are "play", "pause", "stop".
// It returns the the mode and any error encountered.
func SetPlayerMode(hostname string, port int, player_id string, mode string) (string, error) {
	connection_string := fmt.Sprintf("%s:%d", hostname, port)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return "", err
	}
	defer con.Close()

	return setPlayerModeConn(con, player_id, mode)
}

func GetPlayerMode(hostname string, port int, player_id string) (string, error) {
	connection_string := fmt.Sprintf("%s:%d", hostname, port)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return "", err
	}
	defer con.Close()

	return getPlayerModeConn(con, player_id)
}

func TogglePlayerMode(hostname string, port int, player_id string) (string, error) {

	connection_string := fmt.Sprintf("%s:%d", hostname, port)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return "", err
	}
	defer con.Close()

	current_mode, err := getPlayerModeConn(con, player_id)
	if err != nil {
		return "", err
	}

	switch current_mode {
	case "play":
		return setPlayerModeConn(con, player_id, "pause")
	case "pause", "stop":
		return setPlayerModeConn(con, player_id, "play")
	default:
		return current_mode, nil
	}
	return "", nil
}


func getPlayerModeConn(con net.Conn, player_id string) (string, error) {

	cmd := fmt.Sprintf("%s mode ?", player_id);
	s, err := performCommand(con, cmd)
	if err != nil {
		return "", err
	}

	return getTokenFromResponseLineAndDecode(s, 2)
}

func setPlayerModeConn(con net.Conn, player_id string, mode string) (string, error) {

	cmd := fmt.Sprintf("%s mode %s\n", player_id, mode);
	s, err := performCommand(con, cmd)
	if err != nil {
		return "", err
	}

	return getTokenFromResponseLineAndDecode(s, 2)
}