package controller

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/ish-xyz/ykubetest/pkg/client"
	"github.com/ish-xyz/ykubetest/pkg/exporter"
	"github.com/ish-xyz/ykubetest/pkg/loader"
)

const (
	// actions
	APPLY_ACTION  = "apply"
	DELETE_ACTION = "delete"
	EXEC_ACTION   = "exec"
	GET_ACTION    = "get"

	// operators
	REGEX_OPERATOR = "regex:"
)

func runAssertions(code string, resp *client.Response, opsResult *exporter.OperationResult) *exporter.OperationResult {

	env := map[string]interface{}{
		"data": map[string]interface{}{
			"output":  resp.Output,
			"objects": resp.Objects,
			"error":   resp.Error,
		},
	}

	for _, line := range strings.Split(code, ";") {

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Compile expression
		program, err := expr.Compile(line, expr.Env(env))
		opsResult.Status = err == nil
		if err != nil {
			opsResult.AddExpr([2]interface{}{fmt.Sprintf("cannot compile expression: '%s' >> '%v", line, err), opsResult.Status})
			break
		}

		// Run expression
		output, err := expr.Run(program, env)
		opsResult.Status = err == nil
		if err != nil {
			opsResult.AddExpr([2]interface{}{fmt.Sprintf("cannot run expression: '%s' >> '%v", line, err), opsResult.Status})
			break
		}

		opsResult.AddExpr([2]interface{}{line, output})
		opsResult.Status = output != false
		if output == false {
			break
		}
	}

	return opsResult
}

func (ctrl *KubeController) get(ops *loader.TestOperation) *exporter.OperationResult {
	opsResult := &exporter.OperationResult{
		Status:      false,
		ExprResults: make([][2]interface{}, 0),
	}

	resp := ctrl.Client.Get(context.TODO(), ops.ApiVersion, ops.Kind, ops.Namespace, ops.Name, ops.LabelSelector)
	data := runAssertions(ops.Assert, resp, opsResult)

	return data
}

func (ctrl *KubeController) apply(ops *loader.TestOperation) *exporter.OperationResult {

	opsResult := &exporter.OperationResult{
		Status:      false,
		ExprResults: make([][2]interface{}, 0),
	}

	tpl, err := ctrl.Loader.LoadTemplate(ops.Template)
	if err != nil {
		opsResult.AddExpr([2]interface{}{fmt.Sprintf("can't load template %s", ops.Template), false})
		return opsResult
	}
	raw := new(bytes.Buffer)
	templ := template.Must(template.New("").Parse(tpl.Data))
	err = templ.Execute(raw, ops.TemplateValues)
	if err != nil {
		return opsResult
	}

	objects, err := client.GetUnstructuredFromYAML(raw.String())
	if err != nil {
		return opsResult
	}
	resp := ctrl.Client.Apply(context.TODO(), objects)
	opsResult = runAssertions(ops.Assert, resp, opsResult)

	return opsResult
}

func (ctrl *KubeController) delete(ops *loader.TestOperation) *exporter.OperationResult {
	opsResult := &exporter.OperationResult{
		Status:      false,
		ExprResults: make([][2]interface{}, 0),
	}

	tpl, err := ctrl.Loader.LoadTemplate(ops.Template)
	if err != nil {
		opsResult.AddExpr([2]interface{}{fmt.Sprintf("can't load template %s", ops.Template), false})
		return opsResult
	}
	raw := new(bytes.Buffer)
	templ := template.Must(template.New("").Parse(tpl.Data))
	err = templ.Execute(raw, ops.TemplateValues)
	if err != nil {
		return opsResult
	}

	objects, err := client.GetUnstructuredFromYAML(raw.String())
	if err != nil {
		return opsResult
	}
	resp := ctrl.Client.Delete(context.TODO(), objects)
	opsResult = runAssertions(ops.Assert, resp, opsResult)

	return opsResult

}

func (ctrl *KubeController) exec(ops *loader.TestOperation) (string, error) {
	//TODO:
	// run command and get response
	// run assertions
	// return operationResult
	return "", nil
}
