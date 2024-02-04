package helpers

import (
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
