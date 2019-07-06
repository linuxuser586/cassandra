package net

import (
	"errors"
	"net"

	"github.com/linuxuser586/common/pkg/os"
)

var errNoIP = errors.New("no IP for host found")

// PodIP gets the IP address from the POD_IP environment variable.
// If POD_IP is empty, then the first non-loopback IP address is returned.
func PodIP() (string, error) {
	ip := os.Getenv("POD_IP")
	if ip != "" {
		return ip, nil
	}
	addrs, err := InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, a := range addrs {
		if addr, ok := a.(*net.IPNet); ok && !addr.IP.IsLoopback() {
			if addr.IP.To4() != nil {
				return addr.IP.String(), nil
			}
		}
	}
	return "", errNoIP
}
