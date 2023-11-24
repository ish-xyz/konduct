package controller

import (
	"github.com/ish-xyz/ykubetest/pkg/loader"
	"github.com/sirupsen/logrus"
)

func NewController(ldr loader.Loader) Controller {
	return &KubeController{
		Loader: ldr,
		logger: logrus.New().WithField("name", "kube-controller"),
	}
}

func (ctrl *KubeController) Reconcile() error {

	testfiles, err := ctrl.Loader.GetCases()
	if err != nil {
		return err
	}

	for {
		for _, tf := range testfiles {
			tc, err := ctrl.Loader.Load([]string{tf})
			if err != nil {
				ctrl.logger.Warningf("failed to load test file %s", tf)
				continue
			}

			_ = tc

		}
		break
	}
	return nil
}
