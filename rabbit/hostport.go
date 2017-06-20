package rabbit

import (
	"net"
	"strconv"
)

func parseHostPort(hostPort string) (string, uint16, error) {
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
