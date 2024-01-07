package controller

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ControllerManagerOptions struct {
	Config *rest.Config
	Scheme *runtime.Scheme

	MetricsBindAddr string
	MetricsBindPort uint16
	ProbeBindAddr   string

	LeaderElection   bool
	LeaderElectionID string
}

func CreateControllerManager(options *ControllerManagerOptions) (ctrl.Manager, error) {
	return nil, nil
}
