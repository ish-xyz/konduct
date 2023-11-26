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

	data, err := ctrl.Client.Get(context.TODO(), ops.ApiVersion, ops.Kind, ops.Namespace, ops.Name, ops.LabelSelector)

	// TODO keep "not found" in consideration as non error if len is defined
	// if err != nil {
	// 	return "check failed", err
	// }
	fmt.Println(data, err)
	fmt.Println(ops.Assert)
	// for item := range data.Items {
	// 	fmt.Println(item)
	// }
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
