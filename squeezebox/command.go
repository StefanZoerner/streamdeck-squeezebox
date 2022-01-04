package squeezebox

import "net"

func performCommand(connection net.Conn, command string) (string, error) {
	_, err := connection.Write([]byte(command))
	if (err != nil) {
		return "", err
	}

	reply := make([]byte, 2048)
	_, err = connection.Read(reply)
	if (err != nil) {
		return "", err
	}

	return string(reply), nil
}