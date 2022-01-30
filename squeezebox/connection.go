package squeezebox

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type ConnectionProperties struct {
	Hostname string
	HttpPort int
	CliPort int
}

func NewConnectionProperties(hostname string, httpPort, cliPort int) ConnectionProperties {
	return ConnectionProperties{
		Hostname: hostname,
		HttpPort: httpPort,
		CliPort: cliPort,
	}
}

func CheckConnectionCli(cp ConnectionProperties) error {

	connection_string := fmt.Sprintf("%s:%d", cp.Hostname, cp.CliPort)
	con, err := net.Dial("tcp", connection_string)
	if err != nil {
		return err
	}
	defer con.Close()

	result, err := performCommand(con, "version ?")
	if err != nil {
		return err
	} else {
		if ! strings.HasPrefix(result, "version") {
			return errors.New("Unexpected response from server.")
		}
	}

	return nil
}
