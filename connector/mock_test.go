package connector_test

import (
	"net"
	"testing"

	"github.com/filipweidemann/hcloud-kubelet-controller/connector"
	. "github.com/onsi/gomega"
)

var mockConnector connector.MockConnector

func init() {
	mockConnector = connector.MockConnector{}
}

func TestGetServerIPsOutput(t *testing.T) {
	g := NewWithT(t)
	assertedServerIPs := []net.IP{net.ParseIP("192.168.0.100"), net.ParseIP("192.168.0.200")}
	g.Expect(mockConnector.GetServerIPs()).Should(Equal(assertedServerIPs))
}

func TestIsValidForIPsWithCorrectSet(t *testing.T) {
	g := NewWithT(t)
	assertedServerIPs := []net.IP{net.ParseIP("192.168.0.100"), net.ParseIP("192.168.0.200")}
	g.Expect(mockConnector.IsValidForIPs(assertedServerIPs)).Should(BeTrue())
}

func TestIsValidForIPsWithIncorrectSet(t *testing.T) {
	g := NewWithT(t)
	assertedServerIPs := []net.IP{net.ParseIP("192.168.0.250"), net.ParseIP("192.168.0.251")}
	g.Expect(mockConnector.IsValidForIPs(assertedServerIPs)).Should(BeFalse())
}
