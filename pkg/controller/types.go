package controller

import (
	"github.com/ish-xyz/ykubetest/pkg/loader"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	Reconcile() error
}

type KubeController struct {
	Loader loader.Loader
	logger *logrus.Entry
}
