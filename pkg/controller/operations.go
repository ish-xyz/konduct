package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/ish-xyz/ykubetest/pkg/client"
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

func (ctrl *KubeController) get(ops *loader.TestOperation) bool {

	resp := ctrl.Client.Get(context.TODO(), ops.ApiVersion, ops.Kind, ops.Namespace, ops.Name, ops.LabelSelector)
	env := Env{
		Input: resp,
	}

	for _, line := range strings.Split(ops.Assert, ";") {
		line = strings.TrimSpace(line)
		fmt.Println(line)
		if line == "" {
			continue
		}
		program, err := expr.Compile(line, expr.Env(env))
		if err != nil {
			fmt.Println(err)
			return false
		}

		output, err := expr.Run(program, env)
		if err != nil {
			fmt.Println(err)
			return false
		}

		if output == "false" {
			return false
		}
	}

	return true
}

func (ctrl *KubeController) apply(ops *loader.TestOperation) (string, error) {
	fmt.Println("create", ops)
	return "", nil
}

func (ctrl *KubeController) delete(ops *loader.TestOperation) (string, error) {
	return "", nil
}

func (ctrl *KubeController) exec(ops *loader.TestOperation) (string, error) {
	fmt.Println("exec", ops)
	return "", nil
}
