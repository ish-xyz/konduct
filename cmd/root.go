package cmd

import (
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ish-xyz/kubetest/pkg/client"
	"github.com/ish-xyz/kubetest/pkg/controller"
	"github.com/ish-xyz/kubetest/pkg/exporter"
	"github.com/ish-xyz/kubetest/pkg/loader"
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
	err            error
	debug          bool
	controllerMode bool
	interval       int64
	kubeconfig     string
	testsDir       string

	stdout      bool
	pushgateway string

	rootCmd = cobra.Command{
		Long: "A controller and CLI to run e2e tests on Kubernetes",
		Run:  run,
	}
)

func Execute() {
	rootCmd.Execute()
}

func init() {
	// TODO: add label selectors

	// Marked flags
	rootCmd.Flags().StringP("kube-config", "k", "", "path to the kubeconfig file, if empty uses in-cluster method")
	rootCmd.Flags().StringP("dir", "d", "", "Filesystem path to test cases, if empty will load test cases from the Kubernetes API")
	rootCmd.Flags().Int64P("interval", "i", -1, "In controller mode, this settings defines the interval between one test run and the other")
	rootCmd.Flags().BoolP("controller", "c", false, "Run program in controller mode")
	rootCmd.Flags().StringP("pushgateway", "g", "", "Define push gateway address")
	rootCmd.Flags().Bool("stdout", false, "Print report to stdout")

	// Optional flags
	rootCmd.Flags().Bool("debug", false, "Run program in debug mode")

	// MARKERS:

	// can't define --kube-config and --dir together \
	// AND you have to define at least one of them
	rootCmd.MarkFlagsMutuallyExclusive("pushgateway", "stdout")
	rootCmd.MarkFlagsOneRequired("pushgateway", "stdout")

	// can't define --kube-config and --dir together \
	// AND you have to define at least one of them
	rootCmd.MarkFlagsMutuallyExclusive("kube-config", "dir")
	rootCmd.MarkFlagsOneRequired("dir", "kube-config")

	// if --controller is specific --interval should be too
	rootCmd.MarkFlagsRequiredTogether("controller", "interval")

}

func parseFlags(cmd *cobra.Command) error {

	kubeconfig, err = cmd.Flags().GetString("kube-config")
	if err != nil {
		return err
	}

	testsDir, err = cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}

	interval, err = cmd.Flags().GetInt64("interval")
	if err != nil {
		return err
	}

	debug, err = cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}

	controllerMode, err = cmd.Flags().GetBool("controller")
	if err != nil {
		return err
	}

	pushgateway, err = cmd.Flags().GetString("pushgateway")
	if err != nil {
		return err
	}

	stdout, err = cmd.Flags().GetBool("stdout")
	if err != nil {
		return err
	}

	return nil
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

func run(cmd *cobra.Command, args []string) {

	err := parseFlags(cmd)
	checkError(err)

	// Run in debug mode
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	restConfig, err := getRestConfig(expand(kubeconfig))
	checkError(err)

	// Init Kube Client
	clientset, dynclient, err := getClients(restConfig)
	checkError(err)

	kubeclient := client.NewKubeClient(clientset, dynclient, restConfig)

	// Init Tests Loader
	var ldr loader.Loader
	if testsDir != "" {
		// load from filesystem
		ldr, err = loader.NewFSLoader(testsDir, "")
	} else {
		// load from kube api
		ldr, err = loader.NewKubeLoader(kubeclient)
	}
	checkError(err)

	// Init Exporter
	exp := getExporter(debug)

	// Init Controller
	ctrl := controller.NewController(ldr, kubeclient, exp, interval)

	// Execute
	if controllerMode {
		logrus.Fatalln("controller mode is not implemented yet")
	} else {
		err = ctrl.Run(debug)
	}

	checkError(err)
}

func getExporter(verbose bool) exporter.Exporter {
	if stdout {
		return exporter.NewStdoutExporter(verbose)
	}

	// TODO: validate address
	return exporter.NewPrometheusExporter("pushgateway", pushgateway)

}
