package loader

import "github.com/ish-xyz/ykubetest/pkg/client"

type FSloader struct {
	TestsFolder     string
	TemplatesFolder string
}

type KubeLoader struct {
	Client client.Client
}

type Loader interface {
	ListTestCases() ([]string, error)
	LoadTestCase(string) (*TestCase, error)
	LoadTemplate(string) (*TestTemplate, error)
}

// TODO: add Validate() method to struct
type TestTemplate struct {
	Name string `yaml:"name" json:"name"`
	Data string `yaml:"data" json:"data"`
}

// TODO: add Validate() method to struct
type TestCase struct {
	Name string `yaml:"name" json:"name"`

	// +kubebuilder:validation:Required
	Description string `yaml:"description" json:"description"`

	// +kubebuilder:validation:MaxItems=500
	// +kubebuilder:validation:MinItems=1
	Operations []*TestOperation `yaml:"operations" json:"operations"`

	DefaultTimeout string `yaml:"defaultTimeout" json:"defaultTimeout,omitempty"`

	DefaultRetries int `yaml:"defaultRetries" json:"defaultRetries,omitempty"`
}

type TestOperation struct {

	// Used by apply or delete operations only
	Teardown bool `yaml:"teardown" json:"teardown,omitempty"`

	// Used by apply or delete operations only
	Template string `yaml:"template" json:"template,omitempty"`

	// Used by apply or delete operations only
	TemplateValues map[string]bool `yaml:"templateValues" json:"teamplateValues,omitempty"`

	// Used by get operation only
	ApiVersion string `yaml:"apiVersion" json:"apiVersion,omitempty"`

	// Used by get operation only
	Kind string `yaml:"kind" json:"kind,omitempty"`

	// Used by get operation only
	Name string `yaml:"name" json:"name,omitempty"`

	// Used by get operation only
	Namespace string `yaml:"namespace" json:"namespace,omitempty"`

	// Used by get operation only
	LabelSelector string `yaml:"labelSelector" json:"labelSelector,omitempty"`

	// +kubebuilder:validation:Required
	Assert string `yaml:"assert" json:"assert,omitempty"`

	// +kubebuilder:validation:Required
	Action string `yaml:"action" json:"action,omitempty"`

	Retry int `yaml:"retry" json:"retry,omitempty"`

	Interval int `yaml:"interval" json:"interval,omitempty"`

	Wait int `yaml:"wait" json:"wait,omitempty"`
}
