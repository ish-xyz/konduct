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

type TestTemplate struct {
	Name string `yaml:"name" json:"name"`
	Data string `yaml:"data" json:"data"`
}

type TestCase struct {
	Name string `yaml:"name" json:"name"`

	// +kubebuilder:validation:Required
	Description string `yaml:"description" json:"description"`

	// +kubebuilder:validation:MaxItems=500
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	Operations []*TestOperation `yaml:"operations" json:"operations"`

	DefaultTimeout string `yaml:"defaultTimeout" json:"defaultTimeout"`
	DefaultRetries int    `yaml:"defaultRetries" json:"defaultRetries"`
}

type TestOperation struct {

	// Apply/Delete
	Teardown       bool            `yaml:"teardown" json:"teardown"`
	Template       string          `yaml:"template" json:"template"`
	TemplateValues map[string]bool `yaml:"templateValues" json:"teamplateValues"`

	// Get
	ApiVersion    string `yaml:"apiVersion" json:"apiVersion"`
	Kind          string `yaml:"kind" json:"kind"`
	Name          string `yaml:"name" json:"name"`
	Namespace     string `yaml:"namespace" json:"namespace"`
	LabelSelector string `yaml:"labelSelector" json:"labelSelector"`

	// Global
	// +kubebuilder:validation:Required
	Assert string `yaml:"assert" json:"assert"`
	// +kubebuilder:validation:Required
	Action string `yaml:"action" json:"action"`

	// Global Optional
	Retry    int `yaml:"retry" json:"retry"`
	Interval int `yaml:"interval" json:"interval"`
	Wait     int `yaml:"wait" json:"wait"`
}
