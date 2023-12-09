package loader

import (
	"github.com/ish-xyz/kubetest/pkg/client"
)

type FSloader struct {
	TestsFolder     string
	TemplatesFolder string
}

type KubeLoader struct {
	Client client.Client
}

type Loader interface {
	ListTestCases() ([]string, error)
	LoadTestCase(string) (*TestCaseSpec, error)
	LoadTemplate(string) (*TemplateSpec, error)
}

// TODO: add Validate() method to structs: Template and TestCase

// Template definition used to by apply/delete methods
type TemplateSpec struct {
	Name string `yaml:"name" json:"name,omitempty"`
	// +kubebuilder:validation:Required
	Data string `yaml:"data" json:"data"`
}

// TestCase definition used to run e2e tests
type TestCaseSpec struct {

	// +kubebuilder:validation:Required
	Description string `yaml:"description" json:"description"`

	// +kubebuilder:validation:MaxItems=500
	// +kubebuilder:validation:MinItems=1
	Operations []*TestOperation `yaml:"operations" json:"operations"`

	Retry    int `yaml:"retry" json:"retry,omitempty"`
	Interval int `yaml:"interval" json:"interval,omitempty"`
	Wait     int `yaml:"wait" json:"wait,omitempty"`
}

type TestOperation struct {
	Teardown bool   `yaml:"teardown" json:"teardown,omitempty"`
	Template string `yaml:"template" json:"template,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	TemplateValues map[string]interface{} `yaml:"templateValues" json:"templateValues,omitempty"`

	ApiVersion    string `yaml:"apiVersion" json:"apiVersion,omitempty"`
	Kind          string `yaml:"kind" json:"kind,omitempty"`
	LabelSelector string `yaml:"labelSelector" json:"labelSelector,omitempty"`

	Name      string   `yaml:"name" json:"name,omitempty"`
	Namespace string   `yaml:"namespace" json:"namespace,omitempty"`
	Command   []string `yaml:"command" json:"command,omitempty"`

	// +kubebuilder:validation:Required
	Assert string `yaml:"assert" json:"assert,omitempty"`

	// +kubebuilder:validation:Required
	Action   string `yaml:"action" json:"action,omitempty"`
	Retry    int    `yaml:"retry" json:"retry,omitempty"`
	Interval int    `yaml:"interval" json:"interval,omitempty"`
	Wait     int    `yaml:"wait" json:"wait,omitempty"`
}
