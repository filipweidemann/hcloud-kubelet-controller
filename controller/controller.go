package controller

import (
	"context"

	"github.com/filipweidemann/hcloud-kubelet-controller/connector"
	certificatesv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/filipweidemann/hcloud-kubelet-controller/hack/helpers"

	// kubernetes "k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type CertificateSigningRequestReconciler struct {
	Client    client.Client
	Clientset *clientset.Clientset
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

func (r *CertificateSigningRequestReconciler) SetApproval(
	csr *certificatesv1.CertificateSigningRequest,
	approval certificatesv1.RequestConditionType,
) {
	csr.Status.Conditions = append(csr.Status.Conditions, certificatesv1.CertificateSigningRequestCondition{
		Type:           approval,
		Status:         corev1.ConditionTrue,
		Reason:         "automatic hcloud-kubelet-controller validation",
		Message:        "The CSR was approved/denied based on the configured checks inside the hcloud-kubelet-controller",
		LastUpdateTime: metav1.Now(),
	})
}

func (r *CertificateSigningRequestReconciler) UpdateUpstreamResource(
	ctx context.Context,
	req ctrl.Request,
	csr certificatesv1.CertificateSigningRequest,
) (*certificatesv1.CertificateSigningRequest, error) {
	return r.Clientset.CertificatesV1().CertificateSigningRequests().UpdateApproval(
		ctx,
		req.Name,
		&csr,
		metav1.UpdateOptions{},
	)
}

func (r *CertificateSigningRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, rerr error) {
	l := log.FromContext(ctx)
	l.Info("Start reconciliation loop...")

	csr, err := r.FetchCSR(ctx, req)
	if err != nil {
		return res, nil
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

	// Decode CSR to inspect requested content
	x509Request, err := helpers.DecodeCSR(csr.Spec.Request)
	if err != nil {
		l.Info("Error while decoding CSR...")
		return res, rerr
	}

	err = CheckOrganization(*x509Request)
	if err != nil {
		l.Info("Deny approval because of Organization...")
		r.SetApproval(&csr, certificatesv1.CertificateDenied)
		r.UpdateUpstreamResource(ctx, req, csr)
		return ctrl.Result{Requeue: false}, rerr
	}

	// TODO: maybe add CN check here?

	// Upstream Connector IP checks
	ips := x509Request.IPAddresses
	isValid := r.Connector.IsValidForIPs(ips)
	if !isValid {
		l.Info("CSR not valid for IPs: ")
		l.Info("First IP: ")
		l.Info(string(x509Request.IPAddresses[0]))
		r.SetApproval(&csr, certificatesv1.CertificateDenied)
		r.UpdateUpstreamResource(ctx, req, csr)
		return res, rerr
	}

	// Checks are ok, approve & update upstream resource
	l.Info("Setting Approval...")
	r.SetApproval(&csr, certificatesv1.CertificateApproved)
	r.UpdateUpstreamResource(ctx, req, csr)

	return ctrl.Result{
		Requeue: false,
	}, nil
}

func (r *CertificateSigningRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.
		NewControllerManagedBy(mgr).
		For(&certificatesv1.CertificateSigningRequest{}).
		Complete(r)
}
