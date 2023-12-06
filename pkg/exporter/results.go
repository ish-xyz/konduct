package exporter

import (
	"fmt"
)

func NewTestResult(tf string) *TestResult {
	return &TestResult{
		FilePath: tf,
		Name:     "",
		Status:   true,
		Message:  "",
	}
}

func NewOperationResult() *OperationResult {
	return &OperationResult{
		Status:      true,
		Expressions: make([]*ExpressionResult, 0),
	}
}

func (or *OperationResult) AddExpr(exprr *ExpressionResult) {
	or.Expressions = append(or.Expressions, exprr)
}

func (or *OperationResult) Str(printAll bool) string {

	msg := ""
	for _, expr := range or.Expressions {
		if expr.Output == false || printAll {
			msg = fmt.Sprintf("%s\n	* %s >> %v", msg, expr.Expression, expr.Output)
		}
	}
	return msg
}

func (tr *TestResult) Set(status bool, msg string) {
	tr.Message = fmt.Sprintf("%s\n%s", tr.Message, msg)
	if tr.Status && !status {
		tr.Status = status
	}
}
