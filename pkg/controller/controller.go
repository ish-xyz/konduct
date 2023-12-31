package controller

import (
	"fmt"
	"time"

	"github.com/ish-xyz/konduct/pkg/client"
	"github.com/ish-xyz/konduct/pkg/exporter"
	"github.com/ish-xyz/konduct/pkg/loader"
	"github.com/sirupsen/logrus"
)

func NewController(ldr loader.Loader, cl client.Client, exp exporter.Exporter, intrv time.Duration, ro bool) Controller {
	return &KubeController{
		Loader:   ldr,
		Client:   cl,
		Exporter: exp,
		interval: intrv,
		runOnce:  ro,
		logger:   logrus.New().WithField("name", "kube-controller"),
	}
}

func (ctrl *KubeController) singleRun(verbose bool) (*exporter.Report, error) {

	report := exporter.NewReport()

	testfiles, err := ctrl.Loader.ListTestCases()
	if err != nil {
		//TODO: update report as failed
		return nil, err
	}

	for _, tf := range testfiles {

		logrus.Infoln("loading test case:", tf)

		testResult := exporter.NewTestResult(tf)
		testcase, err := ctrl.Loader.LoadTestCase(tf)
		if err != nil {
			ctrl.logger.Errorf("failed to load test file %s %v", tf, err)

			testResult.Set(false, fmt.Sprintf("failed to load testcase in %s", tf))
			report.Add(testResult)
			continue
		}

		// TODO: validate testcase

		for i, op := range testcase.Operations {

			setDefaultTimes(testcase, op)

			opsId := fmt.Sprintf("id-%d", i)

			var opresult = exporter.NewOperationResult()

			if op.Action == loader.GET_ACTION {
				opresult, err = ctrl.get(opsId, op)
			} else if op.Action == loader.APPLY_ACTION {
				opresult, err = ctrl.apply(opsId, op)
			} else if op.Action == loader.DELETE_ACTION {
				opresult, err = ctrl.delete(opsId, op)
			} else if op.Action == loader.EXEC_ACTION {
				opresult, err = ctrl.exec(opsId, op)
			} else {
				err = fmt.Errorf("unkown action '%s'", op.Action)
				opresult.Expressions = append(opresult.Expressions, &exporter.ExpressionResult{Expression: ""})
			}

			if err != nil && len(opresult.Expressions) > 0 {
				opresult.Expressions[len(opresult.Expressions)-1].Expression = fmt.Sprintf(
					"%s %v",
					opresult.Expressions[len(opresult.Expressions)-1].Expression,
					err,
				)
				opresult.Expressions[len(opresult.Expressions)-1].Output = false
			}
			opresult.Status = (err == nil)
			printAll := ctrl.Exporter.IsVerbose() && verbose
			testResult.Set(opresult.Status, opresult.Str(printAll))
		}

		report.Add(testResult)
	}

	return report, nil
}

func setDefaultTimes(tc *loader.TestCaseSpec, op *loader.TestOperation) {
	if op.Interval == 0 && tc.Interval > 0 {
		op.Interval = tc.Interval
	}

	if op.Wait == 0 && tc.Wait > 0 {
		op.Wait = tc.Wait
	}

	if op.Retry == 0 && tc.Retry > 0 {
		op.Retry = tc.Retry
	}
}

func (ctrl *KubeController) Run(verbose bool) error {
	// TODO: add loop for controller
	for {
		report, err := ctrl.singleRun(verbose)
		if err != nil {
			return err
		}

		err = ctrl.Exporter.Export(report)
		if err != nil {
			logrus.Errorln("failed to produce report, error:", err)
		}

		if ctrl.runOnce {
			break
		}

		time.Sleep(time.Second * ctrl.interval)
	}

	return nil
}
