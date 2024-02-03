package controller

import (
	"context"

	certificatesv1 "k8s.io/api/certificates/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	// kubernetes "k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type CertificateSigningRequestReconciler struct {
	Client client.Client
	// cs     *kubernetes.Clientset
	Scheme *runtime.Scheme
}

func (r *CertificateSigningRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, retErr error) {
	l := log.FromContext(ctx)
	l.Info("Start reconciliation loop...")

	csr := certificatesv1.CertificateSigningRequest{}
	err := r.Client.Get(ctx, req.NamespacedName, &csr)

	if apierrors.IsNotFound(err) {
		l.Error(err, "CSR not found.")
		return res, err
	}

	if err != nil {
		l.Error(err, "Error while fetching CSR.")
		return res, err
	}

	l.Info("Found CSR: %v", csr)

	return ctrl.Result{}, nil
}

func (r *CertificateSigningRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.
		NewControllerManagedBy(mgr).
		For(&certificatesv1.CertificateSigningRequest{}).
		Complete(r)
}
