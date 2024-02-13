package controller_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/filipweidemann/hcloud-kubelet-controller/connector"
	"github.com/filipweidemann/hcloud-kubelet-controller/controller"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var testEnv *envtest.Environment
var cfg *rest.Config
var k8sClient client.Client
var k8sClientSet *clientset.Clientset

var testContext context.Context
var testContextCancel context.CancelFunc

func CreateNodeClientSet() *clientset.Clientset {
	// create the node users with group system:masters
	userInfo := envtest.User{Name: "testnode01", Groups: []string{"system:masters"}}
	nodeCfg, err := testEnv.ControlPlane.AddUser(userInfo, cfg)
	if err != nil {
		log.Fatalf("Error during creation of node system user")
	}

	nodeClientSet, err := clientset.NewForConfig(nodeCfg.Config())
	if err != nil {
		log.Fatalf("Error while creating client set for node")
	}

	return nodeClientSet
}

func packageSetup() {
	testContext, testContextCancel = context.WithCancel(context.Background())
	testEnv = &envtest.Environment{}

	cfg, err := testEnv.Start()
	if err != nil {
		log.Fatalf("Failed to setup envtest, %v", err)
	}

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		log.Fatalf("Could not create k8s client: %v", err)
	}
	k8sClientSet = clientset.NewForConfigOrDie(cfg)

	mgrOptions := controller.DefaultControllerManagerOptions(
		scheme.Scheme,
		cfg,
	)
	mgr, err := controller.CreateControllerManager(mgrOptions)
	controller := &controller.CertificateSigningRequestReconciler{
		Client:    mgr.GetClient(),
		Scheme:    mgr.GetScheme(),
		Connector: connector.MockConnector{},
	}

	if err := controller.SetupWithManager(mgr); err != nil {
		log.Fatalf(err.Error(), "unable to create controller", "controller", "Pod")
	}

	go func() {
		err := mgr.Start(testContext)
		if err != nil {
			println("An error occured inside the controller: ", err.Error())
		}
	}()
}

func packageCleanup() {
	testContextCancel()
	testEnv.Stop()
}

func TestMain(m *testing.M) {
	packageSetup()
	gomega.RegisterFailHandler(ginkgo.Fail)
	rc := m.Run()
	packageCleanup()
	os.Exit(rc)
}
