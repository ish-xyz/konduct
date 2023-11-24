package loader

import (
	"github.com/sirupsen/logrus"
)

type FSloader struct {
	Source string
	logger *logrus.Entry
}

type Loader interface {
	Load(files []string) ([]*TestCase, error)
	GetCases() ([]string, error)
}

type TestCase struct {
	Name           string           `yaml:"name"`
	DefaultTimeout string           `yaml:"defaultTimeout"`
	DefaultRetries int              `yaml:"defaultRetries"`
	Operations     []*TestOperation `yaml:"operations"`
}

type TestOperation struct {
	Action         string            `yaml:"action"`
	Template       string            `yaml:"template"`
	TemplateValues map[string]string `yaml:"templateValues"`
}

type ExpectedSpec struct {
	Fail  bool   `yaml:"fail"`
	Error string `yaml:"error"`
	Count int    `yaml:"count"`
}
