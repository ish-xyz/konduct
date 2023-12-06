package cmd

import (
	"fmt"

	"github.com/ish-xyz/ykubetest/pkg/client"
	"github.com/ish-xyz/ykubetest/pkg/controller"
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
	exporter       string
	kubeconfig     string
	sourceDir      string

	rootCmd = cobra.Command{
		Long: "A controller and CLI to run e2e tests on Kubernetes",
		RunE: run,
	}
)

func init() {
	rootCmd.Flags().StringP("kube-config", "k", "~/.kube/config", "path to the kubeconfig file, if empty uses in-cluster method")
	rootCmd.Flags().StringP("source-dir", "s", "", "Filesystem path to test cases, if empty will load tests cases from the Kubernetes API")
	//rootCmd.Flags().StringP("tags", "t", "", "Run tests associated with specific tags")
	rootCmd.Flags().StringP("exporter", "e", "text-file", "Define exporter: prometheus, pushgateway, json-file, text-file")
	rootCmd.Flags().Int64P("interval", "i", 30, "In controller mode, this settings defines the interval between one test run and the other")
	rootCmd.Flags().BoolP("debug", "d", false, "Run program in debug mode")
	rootCmd.Flags().BoolP("controller", "c", false, "Run program in controller mode")

	kubeconfig, err = rootCmd.Flags().GetString("kube-config")
	sourceDir, err = rootCmd.Flags().GetString("source-dir")
	exporter, err = rootCmd.Flags().GetString("exporter")
	interval, err = rootCmd.Flags().GetInt64("interval")
	debug, err = rootCmd.Flags().GetBool("debug")
	controllerMode, err = rootCmd.Flags().GetBool("controller")
}

func Execute() {
	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {

	var restConfig *rest.Config
	var err error
	var templatesDir string
	var report *exporter.Report

	// Run in debug mode
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if sourceDir != "" {
		templatesDir = fmt.Sprintf("%s/%s", sourceDir, "templates")
	}

	// Load kubeconfig
	if kubeconfig == "" {
		restConfig, err = rest.InClusterConfig()
	} else {
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		logrus.Fatalln(err)
	}

	// Allocate clients
	dynclient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		logrus.Fatalln(err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		logrus.Fatalln(err)
	}

	// Start DI
	kubeclient := client.NewKubeClient(clientset, dynclient, restConfig)
	ldr := loader.NewLoader(sourceDir, templatesDir)
	ctrl := controller.NewController(ldr, kubeclient)

	// Execute
	if controllerMode {
		logrus.Warningln("mode not implemented yet")
	} else {
		report, err = ctrl.SingleRun()
	}

	if err != nil {
		logrus.Fatalln(err)
	}

	// Generate report -> TODO: move in controller
	report.Stdout(true)

	return nil
}
