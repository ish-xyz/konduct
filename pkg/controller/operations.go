package controller

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/antonmedv/expr"
	"github.com/ish-xyz/kubetest/pkg/client"
	"github.com/ish-xyz/kubetest/pkg/exporter"
	"github.com/ish-xyz/kubetest/pkg/loader"
	"github.com/sirupsen/logrus"
)

func runAssertions(code string, resp *client.Response) ([]*exporter.ExpressionResult, error) {

	env := map[string]interface{}{
		"data": map[string]interface{}{
			"output":  resp.Output,
			"objects": resp.Objects,
			"error":   resp.Error,
		},
	}

	expressionResults := make([]*exporter.ExpressionResult, 0)

	for _, line := range strings.Split(code, ";") {

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		singleExprResult := &exporter.ExpressionResult{}
		singleExprResult.Expression = line

		expressionResults = append(expressionResults, singleExprResult)

		// Compile expression
		program, err := expr.Compile(line, expr.Env(env))
		if err != nil {
			return expressionResults, err
		}

		// Run expression
		output, err := expr.Run(program, env)
		if err != nil {
			return expressionResults, err
		}

		singleExprResult.Output = output
		if output == false {
			return expressionResults, fmt.Errorf("-") //NOTE: sending an empty error to set the test result to false
		}
	}

	return expressionResults, nil
}

func initialWait(ops *loader.TestOperation) {
	logrus.Infof("waiting %ds to start operation '%s'", ops.Wait, ops.Action)
	time.Sleep(time.Second * time.Duration(ops.Wait))
	return
}

func (ctrl *KubeController) get(opsId string, ops *loader.TestOperation) (*exporter.OperationResult, error) {

	initialWait(ops)

	var err error

	opsResult := &exporter.OperationResult{
		Status:      false,
		Expressions: make([]*exporter.ExpressionResult, 0),
	}

	for ops.Retry >= 0 {

		logrus.Infof("running operation id: %s, action: %s ...\n", opsId, ops.Action)

		resp := ctrl.Client.Get(context.TODO(), ops.ApiVersion, ops.Kind, ops.Namespace, ops.Name, ops.LabelSelector)
		opsResult.Expressions, err = runAssertions(ops.Assert, resp)
		if err == nil {
			break
		}

		ops.Retry--
		time.Sleep(time.Duration(ops.Interval) * time.Second)
	}

	return opsResult, err
}

func (ctrl *KubeController) apply(opsId string, ops *loader.TestOperation) (*exporter.OperationResult, error) {

	initialWait(ops)

	eexp := &exporter.ExpressionResult{Expression: ">> test operation setup"}
	opsResult := &exporter.OperationResult{
		Status:      false,
		Expressions: []*exporter.ExpressionResult{eexp},
	}

	tpl, err := ctrl.Loader.LoadTemplate(ops.Template)
	if err != nil {
		return opsResult, err
	}

	raw := new(bytes.Buffer)
	templ := template.Must(template.New("").Parse(tpl.Data))
	err = templ.Execute(raw, ops.TemplateValues)
	if err != nil {
		return opsResult, err
	}

	objects, err := client.GetUnstructuredFromYAML(raw.String())
	if err != nil {
		return opsResult, err
	}

	for ops.Retry >= 0 {

		logrus.Infof("running operation id: %s, action: %s ...\n", opsId, ops.Action)

		resp := ctrl.Client.Apply(context.TODO(), objects)
		opsResult.Expressions, err = runAssertions(ops.Assert, resp)
		if err == nil {
			break
		}

		ops.Retry--
		time.Sleep(time.Duration(ops.Interval) * time.Second)
	}
	return opsResult, err
}

func (ctrl *KubeController) delete(opsId string, ops *loader.TestOperation) (*exporter.OperationResult, error) {

	initialWait(ops)

	eexp := &exporter.ExpressionResult{Expression: ">> test operation setup"}
	opsResult := &exporter.OperationResult{
		Status:      false,
		Expressions: []*exporter.ExpressionResult{eexp},
	}

	tpl, err := ctrl.Loader.LoadTemplate(ops.Template)
	if err != nil {
		return opsResult, err
	}

	raw := new(bytes.Buffer)
	templ := template.Must(template.New("").Parse(tpl.Data))
	err = templ.Execute(raw, ops.TemplateValues)
	if err != nil {
		return opsResult, err
	}

	objects, err := client.GetUnstructuredFromYAML(raw.String())
	if err != nil {
		return opsResult, err
	}

	for ops.Retry >= 0 {

		logrus.Infof("running operation id: %s, action: %s ...\n", opsId, ops.Action)

		resp := ctrl.Client.Delete(context.TODO(), objects)
		opsResult.Expressions, err = runAssertions(ops.Assert, resp)
		if err == nil {
			break
		}

		ops.Retry--
		time.Sleep(time.Duration(ops.Interval) * time.Second)
	}
	return opsResult, err
}

func (ctrl *KubeController) exec(opsId string, ops *loader.TestOperation) (*exporter.OperationResult, error) {

	initialWait(ops)

	var err error

	opsResult := &exporter.OperationResult{
		Status:      false,
		Expressions: make([]*exporter.ExpressionResult, 0),
	}

	for ops.Retry >= 0 {

		logrus.Infof("running operation id: %s, action: %s ...\n", opsId, ops.Action)

		resp := ctrl.Client.Exec(context.TODO(), ops.Name, ops.Namespace, ops.Command)
		opsResult.Expressions, err = runAssertions(ops.Assert, resp)
		if err == nil {
			break
		}

		ops.Retry--
		time.Sleep(time.Duration(ops.Interval) * time.Second)
	}

	return opsResult, err
}
