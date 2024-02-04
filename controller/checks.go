package controller

import (
	"crypto/x509"
	"errors"
	"strings"
)

func CheckOrganization(req x509.CertificateRequest) error {
	if len(req.Subject.Organization) == 0 {
		return errors.New("Empty CSR Subject Organization, expected system:nodes")
	}

	for _, org := range req.Subject.Organization {
		if !strings.HasPrefix(org, "system:nodes") {
			return errors.New("CSR Subject Organization does not start with system:nodes")
		}
	}

	return nil
}

func CheckCN(req x509.CertificateRequest) error {
	if !strings.HasPrefix(req.Subject.CommonName, "system:node:") {
		return errors.New("CSR Subject CN does not start with system:node:")
	}

	return nil
}
