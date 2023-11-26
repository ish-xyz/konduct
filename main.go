package main

import (
	"github.com/ish-xyz/ykubetest/pkg/client"
	"github.com/ish-xyz/ykubetest/pkg/controller"
	"github.com/ish-xyz/ykubetest/pkg/loader"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	restConfig, _ := clientcmd.BuildConfigFromFlags("", "/home/waffle34/.kube/config")
	dynclient, _ := dynamic.NewForConfig(restConfig)
	clientset, _ := kubernetes.NewForConfig(restConfig)
	kubeclient := client.NewKubeClient(clientset, dynclient, restConfig)
	ldr := loader.NewLoader("filesystem", "./examples")
	ctrl := controller.NewController(ldr, kubeclient)
	ctrl.Run()
}
