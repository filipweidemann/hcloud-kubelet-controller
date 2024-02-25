package connector

import (
	"net"
)

type MockConnector struct{}

func (h MockConnector) GetServerIPs() ([]net.IP, error) {
	// we're testing here. keep it simple, return static list of IPs
	ipLiterals := []string{"192.168.0.100", "192.168.0.200"}
	ips := []net.IP{}
	for _, ip := range ipLiterals {
		ips = append(ips, net.ParseIP(ip))
	}

	return ips, nil
}

func (h MockConnector) IsValidForIPs(ips []net.IP) bool {
	// we know the mock won't error here
	upstreamIPs, _ := h.GetServerIPs()

	// If ANY IP does not match, we will deny the CSR;
	for _, ip := range ips {
		hadMatch := false

		for _, uIP := range upstreamIPs {
			if ip.String() == uIP.String() {
				hadMatch = true
			}
		}

		if !hadMatch {
			return false
		}
	}

	// Otherwise, allow it.
	return true
}
