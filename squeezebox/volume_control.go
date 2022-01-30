package squeezebox

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ChangePlayerVolume(hostname string, cli_port int, player_id string, delta int) (int, error) {

	connection_string := fmt.Sprintf("%s:%d", hostname, cli_port)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return 0, err
	}
	defer con.Close()

	cmd := fmt.Sprintf("%s mixer volume %+d\n", player_id, delta);
	_, err = performCommand(con, cmd)
	if err != nil {
		return 0, err
	}

	cmd = fmt.Sprintf("%s mixer volume ?\n", player_id);
	replyString, err := performCommand(con, cmd)
	if err != nil {
		return 0, err
	}

	parts := strings.Split(replyString, " ")
	current_volume := parts[3]

	volume, err := strconv.Atoi(current_volume)
	if err != nil {
		return 0, err
	}

	return volume, nil
}
