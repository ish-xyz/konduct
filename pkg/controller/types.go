package controller

import (
	"github.com/ish-xyz/ykubetest/pkg/client"
	"github.com/ish-xyz/ykubetest/pkg/loader"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	Run() (*Report, error)
}

type KubeController struct {
	Loader loader.Loader
	Client client.Client
	logger *logrus.Entry
}

type TestResult struct {
	FilePath string
	Name     string
	Status   bool
	Message  string
}

type Report struct {
	Failed   int
	Succeded int
	Status   bool
	Results  []*TestResult
}
