package controller

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

type ControllerManagerOptions struct {
	Config *rest.Config
	Scheme *runtime.Scheme

	MetricsBindAddr string
	MetricsBindPort uint16
	ProbeBindAddr   string

	LeaderElection   bool
	LeaderElectionID string

	K8sConfig *rest.Config
}

func DefaultControllerManagerOptions(s *runtime.Scheme, k8sc *rest.Config) *ControllerManagerOptions {
	return &ControllerManagerOptions{
		Scheme:    s,
		K8sConfig: k8sc,
	}
}

func CreateControllerManager(options *ControllerManagerOptions) (ctrl.Manager, error) {
	if options.K8sConfig == nil {
		println("No kube config specified, using default...")
		options.K8sConfig = ctrl.GetConfigOrDie()
	}

	opts := ctrl.Options{
		Scheme:                  options.Scheme,
		Metrics:                 server.Options{BindAddress: options.MetricsBindAddr},
		HealthProbeBindAddress:  options.ProbeBindAddr,
		LeaderElection:          true,
		LeaderElectionID:        "hcloud-kubelet-controller",
		LeaderElectionNamespace: "kube-system",
	}

	mgr, err := ctrl.NewManager(options.K8sConfig, opts)
	return mgr, err
}
