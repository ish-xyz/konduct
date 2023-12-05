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
		//TODO: update report as failed
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

		for _, op := range testcase.Operations {

			var opresult = exporter.NewOperationResult()

			if op.Action == GET_ACTION {
				opresult, err = ctrl.get(op)
			} else if op.Action == APPLY_ACTION {
				opresult, err = ctrl.apply(op)
				// } else if op.Action == DELETE_ACTION {
				// 	opresult = ctrl.delete(op)
			} else {
				err = fmt.Errorf("unkown action '%s'", op.Action)
				opresult.Expressions = append(opresult.Expressions, &exporter.ExpressionResult{Expression: ""})
			}

			if err != nil && len(opresult.Expressions) > 0 {
				opresult.Expressions[len(opresult.Expressions)-1].Output = fmt.Sprintf("%v", err)
			}
			opresult.Status = (err == nil)
			testResult.Set(opresult.Status, opresult.Str(true))
		}

		report.Add(testResult)
	}

	return report, nil
}
