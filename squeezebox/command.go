package squeezebox

import (
	"errors"
	"fmt"
	"net"
	"net/url"
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

func getTokenFromResponseLineAndDecode(response_line string, n int) (string, error) {
	tokens := strings.Split(response_line, " ")
	if (len(tokens) < n) {
		return "", errors.New(fmt.Sprintf("no token %d in response", n))
	} else {
		decoded, err := url.QueryUnescape(tokens[n])
		if err != nil {
			return "", err
		}

		return decoded, nil
	}
}

func getTagValueFromResponseLine(response_line string, tag_name string) (string, error) {
	value := ""
	var err error = nil

	tokens := strings.Split(response_line, " ")
	for i := 0; i < len(tokens); i++ {
		decoded, _ := url.QueryUnescape(tokens[i])
		if strings.Contains(decoded,":") {
			if strings.HasPrefix(decoded, tag_name + ":") {
				value = decoded[len(tag_name)+1:]
				break
			}
		}
	}

	return value, err
}