package exporter

// Exporters
type PushgatewayExporter struct {
	Address string
}
type FileExporter struct {
	// json or text
	Mode string
	Path string
}

type StdoutExporter struct {
	PrintAll bool
}

// OperationResult functions
type OperationResult struct {
	Status      bool
	Expressions []*ExpressionResult
}

type Exporter interface {
	Export(r *Report) error
	IsVerbose() bool
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
