package squeezebox

import (
	"fmt"
	"net"
	"strconv"
)

func ChangePlayerTrack(hostname string, cliPort int, playerID string, delta int) (int, int, error) {

	var err error
	var trackIndex int
	var trackCount int

	connectionString := fmt.Sprintf("%s:%d", hostname, cliPort)
	con, err := net.Dial("tcp", connectionString)
	if err != nil {
		return 0, 0, err
	}
	defer con.Close()

	// Determine current number of tracks
	//
	cmd := fmt.Sprintf("%s playlist tracks ?", playerID)
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
	trackIndex, err = getTrackIndexConn(con, playerID)
	if err != nil {
		return 0, 0, err
	}

	if delta == 1 {
		if trackCount > trackIndex {
			cmd := fmt.Sprintf("%s playlist index +1", playerID)
			_, err = performCommand(con, cmd)
			if err == nil {
				trackIndex, err = getTrackIndexConn(con, playerID)
			}
		}
	} else if delta == -1 {
		if trackIndex > 1 {
			cmd := fmt.Sprintf("%s playlist index -1", playerID)
			_, err = performCommand(con, cmd)
			if err == nil {
				trackIndex, err = getTrackIndexConn(con, playerID)
			}
		}
	}

	return trackIndex, trackCount, err
}

func getTrackIndexConn(con net.Conn, playerID string) (int, error) {

	var trackIndex = 0
	var err error = nil

	cmd := fmt.Sprintf("%s playlist index ?", playerID)
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
