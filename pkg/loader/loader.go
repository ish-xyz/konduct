package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

const (
	PLACE_HOLDER_EMPTY = "PLACE_HOLDER_EMPTY"
)

func NewLoader(testsFolder, templatesFolder string) Loader {
	return &FSloader{
		TestsFolder:     testsFolder,
		TemplatesFolder: templatesFolder,
	}
}

func (ldr *FSloader) ListTestCases() ([]string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.yaml", ldr.TestsFolder))
	if err != nil {
		return nil, err
	}

	return files, nil
}

// TODO join the following 2 functions with generics

func (ldr *FSloader) LoadTestCase(fname string) (*TestCase, error) {

	var testcase *TestCase

	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &testcase)
	if err != nil {
		return nil, err
	}

	return testcase, nil
}

func (ldr *FSloader) LoadTemplate(fname string) (*TestTemplate, error) {

	var testTempl *TestTemplate

	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &testTempl)
	if err != nil {
		return nil, err
	}

	return testTempl, nil
}
