package sys

import (
	"fmt"
	"net"
)

// GetLocalIPv4 returns the first non-loopback IPv4 address of the host
//
// This function returns the first non-loopback IPv4 address it finds
// in the system. If no such address is found, an error is returned.
func GetLocalIPv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err // failed to get list of interfaces
	}

	// Iterate over all network interfaces
	for _, addr := range addrs {
		// Check if the interface is not a loopback interface
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// Check if the interface is IPv4
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil // found an IPv4 address, return it
			}
		}
	}

	return "", fmt.Errorf("no non-loopback address found") // no IPv4 address found
}
