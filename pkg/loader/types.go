package loader

type FSloader struct {
	TestsFolder     string
	TemplatesFolder string
}

type Loader interface {
	GetTestCases() ([]string, error)
	LoadTestCase(fname string) (*TestCase, error)
	//LoadTemplate()
}

type TestTemplate struct {
	Name string `yaml:"name"`
	Data string `yaml:"data"`
}

type TestCase struct {
	Name           string           `yaml:"name"`
	DefaultTimeout string           `yaml:"defaultTimeout"`
	DefaultRetries int              `yaml:"defaultRetries"`
	Operations     []*TestOperation `yaml:"operations"`
}

type TestOperation struct {
	Action         string                 `yaml:"action"`
	Template       string                 `yaml:"template,omitempty"`
	TemplateValues map[string]interface{} `yaml:"templateValues,omitempty"`
}

type ExpectedSpec struct {
	Fail  bool   `yaml:"fail"`
	Error string `yaml:"error"`
	Count int    `yaml:"count"`
}
