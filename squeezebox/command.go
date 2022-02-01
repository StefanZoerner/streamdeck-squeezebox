package squeezebox

import (
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
	if err != nil {
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
	if strings.Count(response, "\n") == 1 && strings.HasSuffix(response, "\n") {
		response = strings.ReplaceAll(response, "\n", "")
	}

	return response, nil
}

func getTokenFromResponseLineAndDecode(responseLine string, n int) (string, error) {
	tokens := strings.Split(responseLine, " ")
	if len(tokens) < n {
		return "", fmt.Errorf("no token %d in response", n)
	}
	decoded, err := url.QueryUnescape(tokens[n])
	if err != nil {
		return "", err
	}

	return decoded, nil

}

func getTagValueFromResponseLine(responseLine string, tagName string) (string, error) {
	value := ""
	var err error = nil

	tokens := strings.Split(responseLine, " ")
	for i := 0; i < len(tokens); i++ {
		decoded, _ := url.QueryUnescape(tokens[i])
		if strings.Contains(decoded, ":") {
			if strings.HasPrefix(decoded, tagName+":") {
				value = decoded[len(tagName)+1:]
				break
			}
		}
	}

	return value, err
}
