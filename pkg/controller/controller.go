package controller

import (
	"fmt"

	"github.com/ish-xyz/ykubetest/pkg/client"
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

func (ctrl *KubeController) Run() (*Report, error) {

	report := &Report{
		Succeded: 0,
		Failed:   0,
		Status:   true,
		Results:  []*TestResult{},
	}

	testfiles, err := ctrl.Loader.ListTestCases()
	if err != nil {
		return nil, err
	}

	for _, tf := range testfiles {
		testResult := &TestResult{
			FilePath: tf,
			Name:     "",
			Status:   true,
			Message:  "",
		}

		testcase, err := ctrl.Loader.LoadTestCase(tf)
		if err != nil {
			testResult.Status = false
			testResult.Message = fmt.Sprintf("failed to load testcase in %s", tf)
			report.Failed += 1
			report.Status = false
			report.Results = append(report.Results, testResult)

			ctrl.logger.Errorf("failed to load test file %s", tf)
			continue
		}

		testResult.Name = testcase.Name

		for _, ops := range testcase.Operations {

			var err error
			var msg string

			if ops.Action == APPLY_ACTION {
				msg, err = ctrl.apply(ops)
			} else if ops.Action == GET_ACTION {
				fmt.Println(ctrl.get(ops))
			} else if ops.Action == DELETE_ACTION {
				msg, err = ctrl.delete(ops)
			} else if ops.Action == EXEC_ACTION {
				msg, err = ctrl.exec(ops)
			} else {
				err = fmt.Errorf("invalid operation in testcase %s", testcase.Name)
			}

			// report.TotalOperations += 1
			if err != nil {
				testResult.Status = false
				report.Failed += 1
				report.Status = false
			}
			testResult.Message = msg
			report.Results = append(report.Results, testResult)
		}
	}

	return report, nil
}
