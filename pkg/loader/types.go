package loader

type FSloader struct {
	TestsFolder     string
	TemplatesFolder string
}

type Loader interface {
	ListTestCases() ([]string, error)
	LoadTestCase(fname string) (*TestCase, error)
	LoadTemplate(fname string) (*TestTemplate, error)
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

	// Apply & Delete
	Template       string                 `yaml:"template,omitempty"`
	TemplateValues map[string]interface{} `yaml:"templateValues,omitempty"`

	// Get
	ApiVersion    string `yaml:"apiVersion,omitempty"`
	Kind          string `yaml:"kind,omitempty"`
	Name          string `yaml:"name,omitempty"`
	Namespace     string `yaml:"namespace,omitempty"`
	LabelSelector string `yaml:"labelSelector,omitempty"`

	// Global
	Assert *Assert `yaml:"assert"`
	Action string  `yaml:"action"`
}

type Assert struct {
	Fail   bool                   `yaml:"fail,omitempty"`
	Error  string                 `yaml:"error,omitempty" default:"{empty-error}"`
	Count  int                    `yaml:"count,omitempty" default:"-1"`
	Object map[string]interface{} `yaml:"object,omitempty"`
}
