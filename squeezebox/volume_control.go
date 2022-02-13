package squeezebox

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ChangePlayerVolume(cp ConnectionProperties, playerID string, delta int) (int, error) {

	connectionString := fmt.Sprintf("%s:%d", cp.Hostname, cp.CLIPort)
	con, err := net.Dial("tcp", connectionString)
	if err != nil {
		return 0, err
	}
	defer con.Close()

	cmd := fmt.Sprintf("%s mixer volume %+d\n", playerID, delta)
	_, err = performCommand(con, cmd)
	if err != nil {
		return 0, err
	}

	cmd = fmt.Sprintf("%s mixer volume ?\n", playerID)
	replyString, err := performCommand(con, cmd)
	if err != nil {
		return 0, err
	}

	parts := strings.Split(replyString, " ")
	currentVolume := parts[3]

	volume, err := strconv.Atoi(currentVolume)
	if err != nil {
		return 0, err
	}

	return volume, nil
}
