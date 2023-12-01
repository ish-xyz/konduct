package controller

import (
	"github.com/ish-xyz/ykubetest/pkg/client"
	"github.com/ish-xyz/ykubetest/pkg/exporter"
	"github.com/ish-xyz/ykubetest/pkg/loader"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	Exec() (*exporter.Report, error)
}

type KubeController struct {
	Loader loader.Loader
	Client client.Client
	logger *logrus.Entry
}

type Payload struct {
	Message string
	Status  int
}

type Env struct {
	Data  *client.Response
	Print func(format string, a ...any) string
}
