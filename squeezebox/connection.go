package squeezebox

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type ConnectionProperties struct {
	Hostname string
	HTTPPort int
	CLIPort  int
}

func NewConnectionProperties(hostname string, httpPort, cliPort int) ConnectionProperties {
	return ConnectionProperties{
		Hostname: hostname,
		HTTPPort: httpPort,
		CLIPort:  cliPort,
	}
}

func CheckConnectionCLI(cp ConnectionProperties) error {

	connectionString := fmt.Sprintf("%s:%d", cp.Hostname, cp.CLIPort)
	con, err := net.Dial("tcp", connectionString)
	if err != nil {
		return err
	}
	defer con.Close()

	result, err := performCommand(con, "version ?")
	if err != nil {
		return err
	}

	if !strings.HasPrefix(result, "version") {
		return errors.New("unexpected response from server")
	}

	return nil
}

func CheckConnectionToPlayer(hostname string, port int, playerID string) error {

	connectionString := fmt.Sprintf("%s:%d", hostname, port)
	con, err := net.Dial("tcp", connectionString)
	if err != nil {
		return err
	}
	defer con.Close()

	cmd := fmt.Sprintf("%s connected ?\n", playerID)
	_, err = con.Write([]byte(cmd))
	if err != nil {
		return err
	}

	reply := make([]byte, 1024)
	_, err = con.Read(reply)
	if err != nil {
		return err
	}

	sReply := string(reply)
	sReply = strings.ReplaceAll(sReply, "\n", "")
	result := strings.Split(sReply, " ")[2]

	if strings.Contains(result, "%3F") {
		return errors.New("Player " + playerID + " not connected to server.")
	}

	return nil
}
