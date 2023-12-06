package controller

import (
	"github.com/ish-xyz/ykubetest/pkg/client"
	"github.com/ish-xyz/ykubetest/pkg/exporter"
	"github.com/ish-xyz/ykubetest/pkg/loader"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	SingleRun() (*exporter.Report, error)
	Run() error
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
