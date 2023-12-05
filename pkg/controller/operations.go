package controller

import (
	"context"
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
			//fmt.Sprintf("cannot compile expression: '%s' >> '%v", line, err)
			return expressionResults, err
		}

		// Run expression
		output, err := expr.Run(program, env)
		if err != nil {
			//fmt.Sprintf("cannot run expression: '%s' >> '%v", line, err)
			return expressionResults, err
		}

		singleExprResult.Output = output
		if output == false {
			return expressionResults, err
		}
	}

	return expressionResults, nil
}

func (ctrl *KubeController) get(ops *loader.TestOperation) (*exporter.OperationResult, error) {

	var err error

	opsResult := &exporter.OperationResult{
		Status:      false,
		Expressions: make([]*exporter.ExpressionResult, 0),
	}
	resp := ctrl.Client.Get(context.TODO(), ops.ApiVersion, ops.Kind, ops.Namespace, ops.Name, ops.LabelSelector)
	opsResult.Expressions, err = runAssertions(ops.Assert, resp)

	return opsResult, err
}

// func (ctrl *KubeController) apply(ops *loader.TestOperation) *exporter.OperationResult {

// 	opsResult := &exporter.OperationResult{
// 		Status:      false,
// 		Expressions: make([]*exporter.ExpressionResult, 0),
// 	}

// 	tpl, err := ctrl.Loader.LoadTemplate(ops.Template)
// 	if err != nil {
// 		// update expressions, can't load template
// 		return opsResult
// 	}
// 	raw := new(bytes.Buffer)
// 	templ := template.Must(template.New("").Parse(tpl.Data))
// 	err = templ.Execute(raw, ops.TemplateValues)
// 	if err != nil {
// 		// TODO: update expressions
// 		return opsResult
// 	}

// 	objects, err := client.GetUnstructuredFromYAML(raw.String())
// 	if err != nil {
// 		// TODO update expressions
// 		return opsResult
// 	}
// 	resp := ctrl.Client.Apply(context.TODO(), objects)
// 	opsResult.Expressions, opsResult.Status = runAssertions(ops.Assert, resp)

// 	return opsResult
// }

// func (ctrl *KubeController) delete(ops *loader.TestOperation) *exporter.OperationResult {
// 	opsResult := &exporter.OperationResult{
// 		Status:      false,
// 		Expressions: make([]*exporter.ExpressionResult, 0),
// 	}

// 	tpl, err := ctrl.Loader.LoadTemplate(ops.Template)
// 	if err != nil {
// 		// TODO:
// 		return opsResult
// 	}
// 	raw := new(bytes.Buffer)
// 	templ := template.Must(template.New("").Parse(tpl.Data))
// 	err = templ.Execute(raw, ops.TemplateValues)
// 	if err != nil {
// 		// TODO:
// 		return opsResult
// 	}

// 	objects, err := client.GetUnstructuredFromYAML(raw.String())
// 	if err != nil {
// 		// TODO:
// 		return opsResult
// 	}
// 	resp := ctrl.Client.Delete(context.TODO(), objects)
// 	opsResult.Expressions, opsResult.Status = runAssertions(ops.Assert, resp)

// 	return opsResult

// }

// func (ctrl *KubeController) exec(ops *loader.TestOperation) (string, error) {
// 	//TODO:
// 	// run command and get response
// 	// run assertions
// 	// return operationResult
// 	return "", nil
// }
