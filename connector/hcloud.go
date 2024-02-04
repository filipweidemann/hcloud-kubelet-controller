package connector

import (
	"context"
	"net"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type HcloudConnector struct {
	client hcloud.Client
}

func (h HcloudConnector) GetServerIPs() ([]net.IP, error) {
	servers, err := h.client.Server.All(context.Background())
	if err != nil {
		return nil, err
	}

	ips := []net.IP{}
	for _, server := range servers {
		ips = h.getIPFromServer(server, ips)
	}

	return ips, nil
}

func (h HcloudConnector) getIPFromServer(s *hcloud.Server, ips []net.IP) []net.IP {
	ips = append(ips, s.PublicNet.IPv4.IP)
	ips = append(ips, s.PublicNet.IPv6.IP)
	for _, privateNet := range s.PrivateNet {
		ips = append(ips, privateNet.IP)
	}

	return ips
}
