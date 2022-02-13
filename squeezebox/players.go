package squeezebox

import (
	"fmt"
	"net"
	"strconv"
)

type PlayerInfo struct {
	ID   string
	Name string
}

func GetPlayerInfos(cp ConnectionProperties) ([]PlayerInfo, error) {

	connectionString := fmt.Sprintf("%s:%d", cp.Hostname, cp.CLIPort)
	con, err := net.Dial("tcp", connectionString)
	if err != nil {
		return nil, err
	}
	defer con.Close()

	cmd := fmt.Sprintf("player count ?")
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
	for i := 0; i < n; i++ {

		// Get ID for Player i
		cmd := fmt.Sprintf("player id %d ?", i)
		s, err := performCommand(con, cmd)
		if err != nil {
			return nil, err
		}

		playerID, err := getTokenFromResponseLineAndDecode(s, 3)
		if err != nil {
			return nil, err
		}

		cmd = fmt.Sprintf("%s name ?", playerID)
		s, err = performCommand(con, cmd)
		if err != nil {
			return nil, err
		}

		playerName, err := getTokenFromResponseLineAndDecode(s, 2)
		if err != nil {
			return nil, err
		}

		player := PlayerInfo{playerID, playerName}
		infos = append(infos, player)
	}

	return infos, nil
}

func GetPlayerInfo(cp ConnectionProperties, playerID string) (*PlayerInfo, error) {
	connectionString := fmt.Sprintf("%s:%d", cp.Hostname, cp.CLIPort)
	con, err := net.Dial("tcp", connectionString)
	if err != nil {
		return nil, err
	}
	defer con.Close()

	cmd := fmt.Sprintf("%s name ?", playerID)
	s, err := performCommand(con, cmd)
	if err != nil {
		return nil, err
	}

	playerName, err := getTokenFromResponseLineAndDecode(s, 2)
	if err != nil {
		return nil, err
	}

	playerInfo := PlayerInfo{playerID, playerName}

	return &playerInfo, nil
}
