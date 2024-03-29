package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"github.com/filipweidemann/hcloud-kubelet-controller/connector"
	"github.com/filipweidemann/hcloud-kubelet-controller/controller"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var useMockConnector bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&useMockConnector, "use-mock", false,
		"Use mock connector instead of hcloud one.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgrOptions := controller.ControllerManagerOptions{
		Scheme:          scheme,
		ProbeBindAddr:   probeAddr,
		MetricsBindAddr: metricsAddr,
		MetricsBindPort: 9443,
		LeaderElection:  enableLeaderElection,
	}
	mgr, err := controller.CreateControllerManager(&mgrOptions)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	clientset := kubernetes.NewForConfigOrDie(mgrOptions.K8sConfig)
	var conn connector.UpstreamConnector
	if useMockConnector {
		conn = connector.MockConnector{}
	} else {
		token := os.Getenv("HCLOUD_TOKEN")
		if token == "" {
			panic("Please set the environment variable HCLOUD_TOKEN")
		}

		hcloudClient := hcloud.NewClient(hcloud.WithToken(token))
		conn = connector.HcloudConnector{Client: *hcloudClient}
	}

	if err = (&controller.CertificateSigningRequestReconciler{
		Client:    mgr.GetClient(),
		Clientset: clientset,
		Scheme:    mgr.GetScheme(),
		Connector: conn,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CertificateSigningRequest")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
