package controller

import (
	"fmt"

	"github.com/ish-xyz/ykubetest/pkg/client"
	"github.com/ish-xyz/ykubetest/pkg/exporter"
	"github.com/ish-xyz/ykubetest/pkg/loader"
	"github.com/sirupsen/logrus"
)

func NewController(ldr loader.Loader, cl client.Client) Controller {
	return &KubeController{
		Loader: ldr,
		Client: cl,
		logger: logrus.New().WithField("name", "kube-controller"),
	}
}

func (ctrl *KubeController) Exec() (*exporter.Report, error) {

	report := exporter.NewReport()

	testfiles, err := ctrl.Loader.ListTestCases()
	if err != nil {
		// TODO: report should be updated here
		return nil, err
	}

	for _, tf := range testfiles {

		testResult := exporter.NewTestResult(tf)
		testcase, err := ctrl.Loader.LoadTestCase(tf)

		if err != nil {
			ctrl.logger.Errorf("failed to load test file %s", tf)

			testResult.Set(false, fmt.Sprintf("failed to load testcase in %s", tf))
			report.Add(testResult)
			continue
		}

		testResult.Name = testcase.Name

		for _, ops := range testcase.Operations {

			var opsres = exporter.NewOperationResult()

			if ops.Action == GET_ACTION {
				opsres = ctrl.get(ops)
			} else {
				opsres.Status = false
				opsres.AddExpr(
					[2]interface{}{
						fmt.Sprint("invalid operation %s in testcase %s", ops.Action, testcase.Name),
						false,
					})
			}

			testResult.Set(opsres.Status, opsres.Str())
			report.Add(testResult)
		}
	}

	return report, nil
}
