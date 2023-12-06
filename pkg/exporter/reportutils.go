package exporter

import (
	"fmt"
	"strings"
)

func NewTestResult(tf string) *TestResult {
	return &TestResult{
		FilePath: tf,
		Name:     "",
		Status:   true,
		Message:  "",
	}
}

func NewReport() *Report {
	return &Report{
		Succeded: 0,
		Failed:   0,
		Status:   true,
		Results:  []*TestResult{},
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
			msg = fmt.Sprintf("%s\n	* %s %v", msg, expr.Expression, expr.Output)
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

func (r *Report) Add(tr *TestResult) {

	r.Results = append(r.Results, tr)
	if !tr.Status {
		r.Failed++
	} else {
		r.Succeded++
	}

	if !tr.Status && r.Status {
		r.Status = false
	}
}

func (r *Report) Stdout(everything bool) {
	if r.Status {
		fmt.Printf("tests succeded!\ncompleted tests: %d/%d\n", r.Succeded, r.Failed+r.Succeded)
	} else {
		fmt.Printf("tests failed!\ncompleted tests: %d/%d\n", r.Succeded, r.Failed+r.Succeded)
	}

	fmt.Println("results:")

	for _, res := range r.Results {
		fmt.Println(">>", res.FilePath)
		if !res.Status || everything {
			fmt.Printf("	%s\n", strings.TrimSpace(res.Message))
		}
		fmt.Println("--")
	}
}
