package exporter

import "fmt"

// OperationResult functions
type OperationResult struct {
	Status      bool
	ExprResults [][2]interface{}
}

func NewOperationResult() *OperationResult {
	return &OperationResult{
		Status:      true,
		ExprResults: make([][2]interface{}, 0),
	}
}

func (or *OperationResult) AddExpr(expr [2]interface{}) {
	or.ExprResults = append(or.ExprResults, expr)
}

func (or *OperationResult) Set(expr [2]interface{}) {
	or.ExprResults = append(or.ExprResults, expr)
}

func (or *OperationResult) Str(everything bool) string {
	msg := ""
	for _, expr := range or.ExprResults {
		if expr[1] == false || everything {
			msg = fmt.Sprintf("		%s\n%s -> %v", msg, expr[0], expr[1])
		}
	}
	return msg
}

// TestResult functions
type TestResult struct {
	FilePath string
	Name     string
	Status   bool
	Message  string
}

func NewTestResult(tf string) *TestResult {
	return &TestResult{
		FilePath: tf,
		Name:     "",
		Status:   true,
		Message:  "",
	}
}

func (tr *TestResult) Set(status bool, msg string) {
	tr.Message = fmt.Sprintf("%s\n%s", tr.Message, msg)
	if tr.Status && !status {
		tr.Status = status
	}
}

// Report functions
type Report struct {
	Failed   int
	Succeded int
	Status   bool
	Results  []*TestResult
}

func NewReport() *Report {
	return &Report{
		Succeded: 0,
		Failed:   0,
		Status:   true,
		Results:  []*TestResult{},
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
		fmt.Println("tests succeded!")
	} else {
		fmt.Println("tests failed!")
	}
	fmt.Printf("completed tests: %d/%d\n", r.Succeded, r.Failed+r.Succeded)
	fmt.Println("errors:")
	for _, res := range r.Results {

		if !res.Status || everything {
			fmt.Printf("%s\n", res.Message)
		}
	}
}
