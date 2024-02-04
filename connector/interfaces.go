package connector

import "net"

type UpstreamConnector interface {
	GetServerIPs() ([]net.IP, error)
	IsValidForIPs([]net.IP) bool
}
