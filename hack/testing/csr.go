package controller_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"net"

	certificatesv1 "k8s.io/api/certificates/v1"
)

type CSROptions struct {
	NodeName          string
	CN                string
	ExpirationSeconds int32
	DNSName           string
	IPs               []net.IP
}

func CreateTestCSR(csrOpts *CSROptions) certificatesv1.CertificateSigningRequest {
	csr := certificatesv1.CertificateSigningRequest{}
	csr.Name = "test-csr"
	csrOpts.NodeName = "test-node"

	// Assign the well-known Kubernetes signer responsible for the Kubelet CSRs
	// See: https://kubernetes.io/docs/reference/access-authn-authz/certificate-signing-requests/#kubernetes-signers
	csr.Spec.SignerName = certificatesv1.KubeletServingSignerName

	// Set the intended usages and node name
	csr.Spec.Usages = append(csr.Spec.Usages,
		certificatesv1.UsageDigitalSignature,
		certificatesv1.UsageKeyEncipherment,
		certificatesv1.UsageServerAuth,
	)
	csr.Spec.Username = "system:node:" + csrOpts.NodeName
	csrOpts.CN = csr.Spec.Username

	csrOpts.ExpirationSeconds = 600 // value below 600 is not allowed
	csr.Spec.ExpirationSeconds = &csrOpts.ExpirationSeconds

	_, privKey, _ := ed25519.GenerateKey(rand.Reader)

	// Currently, only IP checks are possible with this.
	// Maybe allow settings DNS CNs?
	// It would need to be checked against the IP anyways though...
	x509RequestTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			Organization: []string{"system:nodes"},
			CommonName:   csrOpts.CN,
		},
		IPAddresses: csrOpts.IPs,
	}
	x509RequestTemplate.DNSNames = []string{"test-node.cluster.local"}
	x509Request, _ := x509.CreateCertificateRequest(rand.Reader, &x509RequestTemplate, privKey)

	pemRequest := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: x509Request,
	})

	csr.Spec.Request = pemRequest
	return csr
}
