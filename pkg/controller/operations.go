package controller

import (
	"context"
	"fmt"
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

type Env struct {
	Input *client.Response
	Print func(format string, a ...any) string
}

func runAssertions(code string, resp *client.Response) *exporter.OperationResult {

	opsResult := &exporter.OperationResult{
		Status:      true,
		ExprResults: make([][2]interface{}, 0),
	}

	env := Env{
		Input: resp,
	}

	for _, line := range strings.Split(code, ";") {

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Compile expression
		program, err := expr.Compile(line, expr.Env(env))
		if err != nil {
			opsResult.AddExpr([2]interface{}{fmt.Sprintf("cannot compile expression: '%s'", line), false})
			opsResult.Status = false
			break
		}

		// Run expression
		output, err := expr.Run(program, env)
		if err != nil {
			opsResult.AddExpr([2]interface{}{fmt.Sprintf("cannot run expression: '%s'", line), false})
			opsResult.Status = false
			break
		}

		opsResult.AddExpr([2]interface{}{line, output})
		if output == false {
			opsResult.Status = false
			break
		}
	}

	return opsResult
}

func (ctrl *KubeController) get(ops *loader.TestOperation) *exporter.OperationResult {

	resp := ctrl.Client.Get(context.TODO(), ops.ApiVersion, ops.Kind, ops.Namespace, ops.Name, ops.LabelSelector)
	data := runAssertions(ops.Assert, resp)

	return data
}

func (ctrl *KubeController) apply(ops *loader.TestOperation) (string, error) {
	//TODO:
	// render template
	// run command and get response
	// run assertions
	// return operationResult
	return "", nil
}

func (ctrl *KubeController) delete(ops *loader.TestOperation) (string, error) {
	//TODO:
	// render template
	// run command and get response
	// run assertions
	// return operationResult
	return "", nil
}

func (ctrl *KubeController) exec(ops *loader.TestOperation) (string, error) {
	//TODO:
	// run command and get response
	// run assertions
	// return operationResult
	return "", nil
}
