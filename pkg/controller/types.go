package controller

import (
	"github.com/ish-xyz/kubetest/pkg/client"
	"github.com/ish-xyz/kubetest/pkg/exporter"
	"github.com/ish-xyz/kubetest/pkg/loader"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	Run(bool) error
}

type KubeController struct {
	Loader   loader.Loader
	Client   client.Client
	Exporter exporter.Exporter
	interval int64
	logger   *logrus.Entry
}

type Payload struct {
	Message string
	Status  int
}
