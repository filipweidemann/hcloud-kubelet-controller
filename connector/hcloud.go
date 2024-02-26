package connector

import (
	"context"
	"net"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type HcloudConnector struct {
	Client hcloud.Client
}

func (h HcloudConnector) GetServerIPs() ([]net.IP, error) {
	servers, err := h.Client.Server.All(context.Background())
	println("fetched %v servers!", len(servers))
	if err != nil {
		return nil, err
	}

	ips := []net.IP{}
	for _, server := range servers {
		ips = h.getIPFromServer(server, ips)
	}

	return ips, nil
}

func (h HcloudConnector) IsValidForIPs(csrIPs []net.IP) bool {
	serverIPs, err := h.GetServerIPs()
	if err != nil {
		panic(err)
	}

	validIPs := 0
	for _, serverIP := range serverIPs {
		for _, csrIP := range csrIPs {
			if serverIP.String() == csrIP.String() {
				validIPs += 1
			}
		}
	}

	if len(csrIPs) == validIPs {
		return true
	}

	return false
}

func (h HcloudConnector) getIPFromServer(s *hcloud.Server, ips []net.IP) []net.IP {
	ips = append(ips, s.PublicNet.IPv4.IP)
	ips = append(ips, s.PublicNet.IPv6.IP)
	for _, privateNet := range s.PrivateNet {
		ips = append(ips, privateNet.IP)
	}

	return ips
}
