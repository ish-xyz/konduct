package cmd

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	EXPORTERS_REGEX = "(stdout|pushgateway)"
)

var (
	err error

	rootCmd = cobra.Command{
		Long: "A controller and CLI to run e2e tests on Kubernetes",
		Run:  run,
	}
)

func Execute() {
	rootCmd.Execute()
}

// config:

/*
  kubeconfig: ~/.kube/config
  interval: 4h
  debug: true
  source:
    cluster: false
	filesystem:
	  directory: ./path
  report:
	prometheus:
	  address: 0.0.0.0:9090
	pushgateway:
	  gateway: 0.0.0.0:9091
*/

func init() {
	rootCmd.Flags().StringP("config", "c", "", "pass kubetest config file path")
	rootCmd.MarkFlagRequired("config")
}

func GetRestConfig(kubeconfig string) (*rest.Config, error) {
	var restConfig *rest.Config
	// Load kubeconfig
	if kubeconfig == "" {
		restConfig, err = rest.InClusterConfig()
	} else {
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return restConfig, err
}

func GetClients(restConfig *rest.Config) (*kubernetes.Clientset, *dynamic.DynamicClient, error) {
	// Allocate clients
	dc, errx := dynamic.NewForConfig(restConfig)

	cs, erry := kubernetes.NewForConfig(restConfig)

	if errx != nil {
		return nil, nil, errx
	}

	if erry != nil {
		return nil, nil, erry
	}

	return cs, dc, nil
}

func Expand(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		path = filepath.Join(dir, path[2:])
	}

	return path
}

func checkError(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}

func run(cmd *cobra.Command, args []string) {

	configFile, err := cmd.Flags().GetString("config")
	checkError(err)

	fmt.Println(configFile)
}

// 	// Run in debug mode
// 	if debug {
// 		logrus.SetLevel(logrus.DebugLevel)
// 	}

// 	restConfig, err := getRestConfig(expand(kubeconfig))
// 	checkError(err)

// 	// Init Kube Client
// 	clientset, dynclient, err := getClients(restConfig)
// 	checkError(err)

// 	kubeclient := client.NewKubeClient(clientset, dynclient, restConfig)

// 	// Init Tests Loader
// 	var ldr loader.Loader
// 	if testsDir != "" {
// 		// load from filesystem
// 		ldr, err = loader.NewFSLoader(testsDir, "")
// 	} else {
// 		// load from kube api
// 		ldr, err = loader.NewKubeLoader(kubeclient)
// 	}
// 	checkError(err)

// 	// Init Exporter
// 	exp := getExporter(debug)

// 	// Init Controller
// 	ctrl := controller.NewController(ldr, kubeclient, exp, interval)

// 	// Execute
// 	if controllerMode {
// 		logrus.Fatalln("controller mode is not implemented yet")
// 	} else {
// 		err = ctrl.Run(debug)
// 	}

// 	checkError(err)
// }

// func getExporter(verbose bool) exporter.Exporter {
// 	if stdout {
// 		return exporter.NewStdoutExporter(verbose)
// 	}

// 	// TODO: validate address
// 	return exporter.NewPrometheusExporter("pushgateway", pushgateway)

// }
