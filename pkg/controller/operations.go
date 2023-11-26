package controller

import (
	"context"
	"fmt"

	"github.com/ish-xyz/ykubetest/pkg/loader"
)

const (
	APPLY_OPERATION  = "apply"
	DELETE_OPERATION = "delete"
	EXEC_OPERATION   = "exec"
	GET_OPERATION    = "get"
)

func (ctrl *KubeController) get(ops *loader.TestOperation) (string, error) {

	data, _ := ctrl.Client.Get(context.TODO(), ops.ApiVersion, ops.Kind, ops.Namespace, ops.LabelSelector)

	if ops.Name != "" {
		for _, d := range data.Items {
			if d.GetName() == ops.Name {
				fmt.Println(d)
			}
		}
	}
	return "", nil

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
