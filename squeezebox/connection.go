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
