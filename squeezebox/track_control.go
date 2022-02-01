package squeezebox

import (
	"fmt"
	"net"
	"strconv"
)

func ChangePlayerTrack(hostname string, cli_port int, player_id string, delta int) (int, int, error) {

	var err error
	var trackIndex int
	var trackCount int

	connection_string := fmt.Sprintf("%s:%d", hostname, cli_port)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return 0, 0, err
	}
	defer con.Close()

	// Determine current number of tracks
	//
	cmd := fmt.Sprintf("%s playlist tracks ?", player_id)
	resp, err := performCommand(con, cmd)
	if err == nil {
		s, err := getTokenFromResponseLineAndDecode(resp, 3)
		if err == nil {
			trackCount, err = strconv.Atoi(s)
		}
	}
	if err != nil {
		return 0, 0, err
	}

	if trackCount == 0 {
		return 0, 0, nil
	}

	// Determine current track number
	trackIndex, err = getTrackIndexConn(con, player_id)
	if err != nil {
		return 0, 0, err
	}

	if delta == 1 {
		if trackCount > trackIndex {
			cmd := fmt.Sprintf("%s playlist index +1", player_id)
			_, err = performCommand(con, cmd)
			if err == nil {
				trackIndex, err = getTrackIndexConn(con, player_id)
			}
		}
	} else if delta == -1 {
		if trackIndex > 1 {
			cmd := fmt.Sprintf("%s playlist index -1", player_id)
			_, err = performCommand(con, cmd)
			if err == nil {
				trackIndex, err = getTrackIndexConn(con, player_id)
			}
		}
	}

	return trackIndex, trackCount, err
}

func getTrackIndexConn(con net.Conn, player_id string) (int, error) {

	var trackIndex = 0
	var err error = nil

	cmd := fmt.Sprintf("%s playlist index ?", player_id)
	resp, err := performCommand(con, cmd)
	if err == nil {
		s, err := getTokenFromResponseLineAndDecode(resp, 3)
		if err == nil {
			trackIndex, err = strconv.Atoi(s)
			if err == nil {
				trackIndex = trackIndex + 1
			}
		}
	}

	return trackIndex, err
}
