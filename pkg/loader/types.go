package loader

type FSloader struct {
	Source string
}

type Loader interface {
	GetTestCases() ([]string, error)
	LoadTestCase(fname string) (*TestCase, error)
	//LoadTemplate()
}

type TestCase struct {
	Name           string           `yaml:"name"`
	DefaultTimeout string           `yaml:"defaultTimeout"`
	DefaultRetries int              `yaml:"defaultRetries"`
	Operations     []*TestOperation `yaml:"operations"`
}

type TestOperation struct {
	Action         string            `yaml:"action"`
	Template       string            `yaml:"template,omitempty"`
	TemplateValues map[string]string `yaml:"templateValues,omitempty"`
}

type ExpectedSpec struct {
	Fail  bool   `yaml:"fail"`
	Error string `yaml:"error"`
	Count int    `yaml:"count"`
}
