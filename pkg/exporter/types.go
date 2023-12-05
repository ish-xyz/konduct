package exporter

// OperationResult functions
type OperationResult struct {
	Status      bool
	Expressions []*ExpressionResult
}

type ExpressionResult struct {
	Expression string
	Output     interface{}
}

// TestResult functions
type TestResult struct {
	FilePath string
	Name     string
	Status   bool
	Message  string
}

// Report functions
type Report struct {
	Failed   int
	Succeded int
	Status   bool
	Results  []*TestResult
}