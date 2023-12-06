package main

import (
	"github.com/ish-xyz/ykubetest/pkg/client"
	"github.com/ish-xyz/ykubetest/pkg/controller"
	"github.com/ish-xyz/ykubetest/pkg/loader"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	restConfig, err := clientcmd.BuildConfigFromFlags("", "/Users/ishamaraia/.kube/config")
	if err != nil {
		logrus.Fatalln(err)
	}
	dynclient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		logrus.Fatalln(err)
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		logrus.Fatalln(err)
	}
	kubeclient := client.NewKubeClient(clientset, dynclient, restConfig)
	ldr := loader.NewLoader("./examples", "./examples/templates")
	ctrl := controller.NewController(ldr, kubeclient)
	report, err := ctrl.SingleRun()

	if err != nil {
		panic(err)
	}

	report.Stdout(true)

	return
}
