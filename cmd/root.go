package cmd

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ish-xyz/ykubetest/pkg/client"
	"github.com/ish-xyz/ykubetest/pkg/controller"
	"github.com/ish-xyz/ykubetest/pkg/exporter"
	"github.com/ish-xyz/ykubetest/pkg/loader"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	err            error
	debug          bool
	controllerMode bool
	interval       int64
	//exporterType   string
	kubeconfig   string
	sourceDir    string
	templatesDir string

	rootCmd = cobra.Command{
		Long: "A controller and CLI to run e2e tests on Kubernetes",
		Run:  run,
	}
)

func init() {

	rootCmd.Flags().StringP("kube-config", "k", "~/.kube/config", "path to the kubeconfig file, if empty uses in-cluster method")
	rootCmd.Flags().StringP("source-dir", "s", "", "Filesystem path to test cases, if empty will load tests cases from the Kubernetes API")
	//rootCmd.Flags().StringP("labels", "l", "", "Run tests associated with specific labels")
	//rootCmd.Flags().StringP("exporter", "e", "stdout", "Define exporter: stdout, prometheus, pushgateway, json-file, text-file")
	rootCmd.Flags().Int64P("interval", "i", 30, "In controller mode, this settings defines the interval between one test run and the other")
	rootCmd.Flags().BoolP("debug", "d", false, "Run program in debug mode")
	rootCmd.Flags().BoolP("controller", "c", false, "Run program in controller mode")

	kubeconfig, err = rootCmd.Flags().GetString("kube-config")

	checkError(err)

	sourceDir, err = rootCmd.Flags().GetString("source-dir")
	checkError(err)

	// exporterType, err = rootCmd.Flags().GetString("exporter")
	// checkError(err)

	interval, err = rootCmd.Flags().GetInt64("interval")
	checkError(err)

	debug, err = rootCmd.Flags().GetBool("debug")
	checkError(err)

	controllerMode, err = rootCmd.Flags().GetBool("controller")
	checkError(err)

}
func expand(path string) string {
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

func Execute() {
	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {

	var err error

	// Run in debug mode
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	kubeconfig = expand(kubeconfig)

	restConfig, err := getRestConfig(kubeconfig)
	checkError(err)

	// Init Kube Client
	clientset, dynclient, err := getClients(restConfig)
	checkError(err)

	kubeclient := client.NewKubeClient(clientset, dynclient, restConfig)

	// Init Tests Loader
	var ldr loader.Loader
	if sourceDir != "" {
		// load from filesystem
		ldr, err = loader.NewFSLoader(sourceDir, templatesDir)
	} else {
		// load from kube api
		logrus.Fatalln("CRDs are not implemented yet")
		fmt.Println("here")
		ldr, err = loader.NewKubeLoader(kubeclient)
	}
	checkError(err)

	// Init Exporter
	exp := exporter.NewStdoutExporter(debug)

	// Init Controller
	ctrl := controller.NewController(ldr, kubeclient, exp, interval)

	// Execute
	if controllerMode {
		logrus.Fatalln("controller mode is not implemented yet")
	} else {
		err = ctrl.Run()
	}

	checkError(err)
}

func getRestConfig(kubeconfig string) (*rest.Config, error) {
	var restConfig *rest.Config
	// Load kubeconfig
	if kubeconfig == "" {
		restConfig, err = rest.InClusterConfig()
	} else {
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return restConfig, err
}

func getClients(restConfig *rest.Config) (*kubernetes.Clientset, *dynamic.DynamicClient, error) {
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
