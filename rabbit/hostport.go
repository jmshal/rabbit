package rabbit

import (
	"net"
	"strconv"
	"strings"
)

func parseHostPort(hostPort string) (string, uint16, error) {
	if strings.Index(hostPort, ":") == -1 {
		return hostPort, 0, nil
	}

	host, portString, err := net.SplitHostPort(hostPort)

	if err != nil {
		return "", 0, err
	}

	if portString != "" {
		port, err := strconv.Atoi(portString)

		if err != nil {
			return "", 0, err
		}

		return host, uint16(port), nil
	}

	return host, 0, nil
}
