package connector

import "net"

type UpstreamConnector interface {
	GetServerIPs() []net.IP
	IsValidForIPs([]net.IP) bool
}
