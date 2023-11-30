package exporter

import "fmt"

type OperationResult struct {
	Status      bool
	ExprResults [][2]interface{}
}

type TestResult struct {
	FilePath string
	Name     string
	Status   bool
	Message  string
}

type Report struct {
	Failed   int
	Succeded int
	Status   bool
	Results  []*TestResult
}

func NewOperationResult() *OperationResult {
	return &OperationResult{
		Status:      true,
		ExprResults: make([][2]interface{}, 0),
	}
}

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

func (tr *TestResult) Set(status bool, msg string) {
	tr.Message = msg
	tr.Status = status
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

func (or *OperationResult) AddExpr(expr [2]interface{}) {
	or.ExprResults = append(or.ExprResults, expr)
}

func (or *OperationResult) Str() string {
	msg := "failed operations:"
	for _, expr := range or.ExprResults {
		if expr[1] != true {
			msg = fmt.Sprintf("%s\n%s -> %s", msg, expr[0], expr[1])
		}
	}
	return msg
}
