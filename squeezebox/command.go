package squeezebox

import (
	"net"
	"strings"
)

func performCommand(connection net.Conn, command string) (string, error) {

	// Append newline character to command if missing
	if !strings.HasSuffix(command, "\n") {
		command += "\n"
	}

	// Send to server
	_, err := connection.Write([]byte(command))
	if (err != nil) {
		return "", err
	}

	// Receive server reply in byte array
	buffSize := 2048
	reply := make([]byte, buffSize)
	n, err := connection.Read(reply)
	if err != nil {
		return "", err
	}

	// Convert byte array to String
	response := string(reply[:n])

	// Remove trailing newline character, if response is one line
	if strings.Count(response, "\n") == 1 &&  strings.HasSuffix(response, "\n") {
		response = strings.ReplaceAll(response, "\n", "")
	}

	return response, nil
}