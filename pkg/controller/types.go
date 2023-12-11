package controller

import (
	"time"

	"github.com/ish-xyz/konduct/pkg/client"
	"github.com/ish-xyz/konduct/pkg/exporter"
	"github.com/ish-xyz/konduct/pkg/loader"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	Run(bool) error
}

type KubeController struct {
	Loader   loader.Loader
	Client   client.Client
	Exporter exporter.Exporter
	interval time.Duration
	logger   *logrus.Entry
	runOnce  bool
}

type Payload struct {
	Message string
	Status  int
}
