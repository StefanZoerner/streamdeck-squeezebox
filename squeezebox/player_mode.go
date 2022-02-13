package squeezebox

import (
	"fmt"
	"net"
)

// SetPlayerMode sets the mode of the player with playerID to mode.
// Possible values are "play", "pause", "stop".
// It returns the the mode and any error encountered.
func SetPlayerMode(cp ConnectionProperties, playerID string, mode string) (string, error) {
	connectionString := fmt.Sprintf("%s:%d", cp.Hostname, cp.CLIPort)
	con, err := net.Dial("tcp", connectionString)
	if err != nil {
		return "", err
	}
	defer con.Close()

	return setPlayerModeConn(con, playerID, mode)
}

// GetPlayerMode retrieves the mode of the player with playerID.
// Possible values are "play", "pause", "stop".
// It returns the the mode and any error encountered.
func GetPlayerMode(cp ConnectionProperties, playerID string) (string, error) {
	connectionString := fmt.Sprintf("%s:%d", cp.Hostname, cp.CLIPort)
	con, err := net.Dial("tcp", connectionString)
	if err != nil {
		return "", err
	}
	defer con.Close()

	return getPlayerModeConn(con, playerID)
}

// TogglePlayerMode changes the mode of the player with playerID.
// From play to pause, from stop or pause to play.
// It returns the new mode and any error encountered.
func TogglePlayerMode(cp ConnectionProperties, playerID string) (string, error) {

	connectionString := fmt.Sprintf("%s:%d", cp.Hostname, cp.CLIPort)
	con, err := net.Dial("tcp", connectionString)
	if err != nil {
		return "", err
	}
	defer con.Close()

	currentMode, err := getPlayerModeConn(con, playerID)
	if err != nil {
		return "", err
	}

	switch currentMode {
	case "play":
		return setPlayerModeConn(con, playerID, "pause")
	case "pause", "stop":
		return setPlayerModeConn(con, playerID, "play")
	default:
		return currentMode, nil
	}
}

func getPlayerModeConn(con net.Conn, playerID string) (string, error) {

	cmd := fmt.Sprintf("%s mode ?", playerID)
	s, err := performCommand(con, cmd)
	if err != nil {
		return "", err
	}

	return getTokenFromResponseLineAndDecode(s, 2)
}

func setPlayerModeConn(con net.Conn, playerID string, mode string) (string, error) {

	cmd := fmt.Sprintf("%s mode %s\n", playerID, mode)
	s, err := performCommand(con, cmd)
	if err != nil {
		return "", err
	}

	return getTokenFromResponseLineAndDecode(s, 2)
}
