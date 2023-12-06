package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

const (
	TEMPLATE_FOLDER = "templates"
)

func NewFSLoader(testsFolder, templatesFolder string) Loader {
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

func (ldr *FSloader) LoadTemplate(tname string) (*TestTemplate, error) {

	var testTempl *TestTemplate

	tplname := fmt.Sprintf("%s/%s/%s.yaml", ldr.TestsFolder, TEMPLATE_FOLDER, tname)

	data, err := os.ReadFile(tplname)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &testTempl)
	if err != nil {
		return nil, err
	}

	return testTempl, nil
}
