package loader

import (
	"fmt"
)

const (
	GET_ACTION    = "get"
	APPLY_ACTION  = "apply"
	EXEC_ACTION   = "exec"
	DELETE_ACTION = "delete"
)

func (tc *TestCaseSpec) Validate() error {

	err := validateOps(tc.Operations)
	if err != nil {
		return err
	}

	return nil
}

func validateOps(ops []*TestOperation) error {
	for _, op := range ops {
		// Assert should always be defined
		if op.Assert == "" {
			return fmt.Errorf("no assert specified")
		}

		if op.Action == GET_ACTION {
			err := validGet(op)
			if err != nil {
				return err
			}
		} else if op.Action == APPLY_ACTION {
			err := validApplyDelete(op)
			if err != nil {
				return err
			}
		} else if op.Action == DELETE_ACTION {
			err := validApplyDelete(op)
			if err != nil {
				return err
			}
		} else if op.Action == EXEC_ACTION {
			err := validExec(op)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func validGet(op *TestOperation) error {
	if op.Kind == "" || op.ApiVersion == "" {
		return fmt.Errorf("no kind or apiVersion specified")
	}
	if op.Name == "" && op.LabelSelector == "" {
		return fmt.Errorf("no name or labelSelector specified")
	}

	return nil
}

func validApplyDelete(op *TestOperation) error {
	if op.Template == "" {
		return fmt.Errorf("no template defined")
	}

	return nil
}

func validExec(op *TestOperation) error {
	if len(op.Command) == 0 {
		return fmt.Errorf("no command defined for exec operation")
	}

	if op.Name == "" || op.Namespace == "" {
		return fmt.Errorf("no namespace or name defined for exec operation")
	}

	return nil
}
