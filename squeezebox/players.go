package squeezebox

import (
	"fmt"
	"net"
	"strconv"
)

type PlayerInfo struct {
	Id string
	Name string
}

func GetPlayers(hostname string, port int) ([]PlayerInfo, error) {

	connection_string := fmt.Sprintf("%s:%d", hostname, port)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return nil, err
	}
	defer con.Close()

	cmd := fmt.Sprintf("player count ?");
	s, err := performCommand(con, cmd)
	if err != nil {
		return nil, err
	}

	count, err := getTokenFromResponseLineAndDecode(s, 2)
	if err != nil {
		return nil, err
	}
	n, err := strconv.Atoi(count)
	if err != nil {
		return nil, err
	}

	infos := []PlayerInfo{}
	for i := 0; i < n;i++ {

		// Get Id for Player i
		cmd := fmt.Sprintf("player id %d ?", i);
		s, err := performCommand(con, cmd)
		if err != nil {
			return nil, err
		}

		player_id, err := getTokenFromResponseLineAndDecode(s, 3)
		if err != nil {
			return nil, err
		}

		cmd = fmt.Sprintf("%s name ?", player_id);
		s, err = performCommand(con, cmd)
		if err != nil {
			return nil, err
		}

		player_name, err := getTokenFromResponseLineAndDecode(s, 2)
		if err != nil {
			return nil, err
		}

		player := PlayerInfo{ player_id, player_name }
		infos = append(infos, player)
	}

	return infos, nil
}

func GetPlayerInfo(hostname string, port int, player_id string) (*PlayerInfo, error) {
	connection_string := fmt.Sprintf("%s:%d", hostname, port)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return nil, err
	}
	defer con.Close()

	cmd := fmt.Sprintf("%s name ?", player_id);
	s, err := performCommand(con, cmd)
	if err != nil {
		return nil, err
	}

	player_name, err := getTokenFromResponseLineAndDecode(s, 2)
	if err != nil {
		return nil, err
	}

	playerInfo := PlayerInfo{ player_id, player_name }

	return &playerInfo, nil
}