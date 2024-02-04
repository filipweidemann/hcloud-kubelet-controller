package helpers

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	certificatesv1 "k8s.io/api/certificates/v1"
)

func GetCSRApproval(csr *certificatesv1.CertificateSigningRequest) (approved bool, denied bool) {
	for _, condition := range csr.Status.Conditions {
		if condition.Type == certificatesv1.CertificateApproved {
			approved = true
		}

		if condition.Type == certificatesv1.CertificateDenied {
			denied = true
		}
	}

	return approved, denied
}

func DecodeCSR(pemBytes []byte) (*x509.CertificateRequest, error) {
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil || pemBlock.Type != "CERTIFICATE REQUEST" {
		return nil, errors.New("Not a Certificate Request")
	}

	csr, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	return csr, err
}
