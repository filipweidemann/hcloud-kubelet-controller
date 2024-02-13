package controller_test

import (
	"net"
	"testing"
	"time"

	"github.com/filipweidemann/hcloud-kubelet-controller/hack/testing"
	"github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestApprovalForValidCSR(t *testing.T) {
	g := gomega.NewWithT(t)

	csrCfg := controller_test.CSROptions{
		IPs:      []net.IP{net.ParseIP("192.168.0.100")},
		CN:       "192.168.0.100",
		NodeName: "testnode01",
		DNSName:  "testnode01.company.tld",
	}
	csr := controller_test.CreateTestCSR(&csrCfg)

	nodeClientSet := CreateNodeClientSet()
	_, err := nodeClientSet.CertificatesV1().CertificateSigningRequests().Create(testContext, &csr, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("Error during creation of CSR: %v", err)
	}

	g.Eventually(func() bool {
		csr, err := k8sClientSet.CertificatesV1().CertificateSigningRequests().Get(testContext, "test-csr", metav1.GetOptions{})

		if err != nil {
			t.Error("Could not fetch updated CSR")
		}

		if csr.Status.Certificate != nil {
			return true
		}

		return false
	}, time.Second*2, time.Millisecond*500).Should(gomega.BeTrue())
}
