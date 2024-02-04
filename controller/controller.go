package controller

import (
	"context"

	"github.com/filipweidemann/hcloud-kubelet-controller/connector"
	certificatesv1 "k8s.io/api/certificates/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/filipweidemann/hcloud-kubelet-controller/hack/helpers"

	// kubernetes "k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type CertificateSigningRequestReconciler struct {
	Client    client.Client
	Scheme    *runtime.Scheme
	Connector connector.UpstreamConnector
}

func (r *CertificateSigningRequestReconciler) FetchCSR(ctx context.Context, req ctrl.Request) (csr certificatesv1.CertificateSigningRequest, err error) {
	l := log.FromContext(ctx)
	csr = certificatesv1.CertificateSigningRequest{}
	err = r.Client.Get(ctx, req.NamespacedName, &csr)

	if apierrors.IsNotFound(err) {
		l.Info("CSR not found.")
		return csr, err
	}

	if err != nil {
		l.Error(err, "Error while fetching CSR.")
		return csr, err
	}

	l.Info("Found CSR: ")
	l.Info(csr.Name)

	return csr, err
}

func (r *CertificateSigningRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, rerr error) {
	l := log.FromContext(ctx)
	l.Info("Start reconciliation loop...")

	csr, err := r.FetchCSR(ctx, req)
	if err != nil {
		return res, err
	}

	// ignore CSRs that are not meant to be processed by the KubeletServingSigner
	if csr.Spec.SignerName != certificatesv1.KubeletServingSignerName {
		l.Info("Ignore CSR not meant to be processed by controller...")
		return res, rerr
	}

	approved, denied := helpers.GetCSRApproval(&csr)
	if approved || denied {
		l.Info("This CSR was already processed, ignoring...")
		return res, rerr
	}

	return res, err
}

func (r *CertificateSigningRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.
		NewControllerManagedBy(mgr).
		For(&certificatesv1.CertificateSigningRequest{}).
		Complete(r)
}
