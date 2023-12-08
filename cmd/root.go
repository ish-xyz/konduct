package cmd

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/ish-xyz/kubetest/pkg/client"
	"github.com/ish-xyz/kubetest/pkg/controller"
	"github.com/ish-xyz/kubetest/pkg/exporter"
	"github.com/ish-xyz/kubetest/pkg/loader"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	SOURCE_CLUSTER = "cluster"
	SOURCE_FS      = "filesystem"

	EXPORTER_STDOUT      = "stdout"
	EXPORTER_PUSHGATEWAY = "pushgateway"
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

func init() {
	rootCmd.Flags().StringP("config", "c", "", "pass kubetest config file path")
	rootCmd.MarkFlagRequired("config")
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

	configFile, err := cmd.Flags().GetString("config")
	checkError(err)

	lviper := viper.New()

	lviper.SetConfigFile(configFile)
	err = lviper.ReadInConfig()
	checkError(err)

	kcfg := lviper.GetString("kubeconfig")
	if kcfg == "" {
		checkError(fmt.Errorf("empty kubeconfig path in config file"))
	}

	// Run in debug mode
	if lviper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Init Clients
	restConfig, err := getRestConfig(expand(kcfg))
	checkError(err)

	clientset, dynclient, err := getClients(restConfig)
	checkError(err)

	kubeclient := client.NewKubeClient(clientset, dynclient, restConfig)

	// Init loader
	var ldr loader.Loader
	sourceType := lviper.GetString("source.type")

	if sourceType == SOURCE_CLUSTER {
		ldr, err = loader.NewKubeLoader(kubeclient)
	} else if sourceType == SOURCE_FS {
		ldr, err = loader.NewFSLoader(lviper.GetString("source.directory"), "")
	} else {
		err = fmt.Errorf("invalid source.type in config file")
	}
	checkError(err)

	// Init exporter
	var exp exporter.Exporter
	exporterType := lviper.GetString("report.type")
	if exporterType == EXPORTER_PUSHGATEWAY {
		exp, err = exporter.NewPushgatewayExporter(lviper.GetString("report.address"))
	} else if exporterType == EXPORTER_STDOUT {
		exp = exporter.NewStdoutExporter(lviper.GetBool("debug"))
	} else {
		err = fmt.Errorf("invalid report.type in config file")
	}
	checkError(err)

	// Init Controller
	duration, err := time.ParseDuration(lviper.GetString("interval"))
	if err != nil {
		duration = time.Duration(time.Hour * 4)
	}
	ctrl := controller.NewController(ldr, kubeclient, exp, duration, lviper.GetBool("runOnce"))

	checkError(ctrl.Run(lviper.GetBool("debug")))
}
