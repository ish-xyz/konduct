package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

const (
	TEMPLATE_FOLDER = "templates"
)

func NewFSLoader(testsFolder, templatesFolder string) (Loader, error) {

	if testsFolder == "" {
		return nil, fmt.Errorf("you must define a source folder for your tests")
	}
	if templatesFolder == "" {
		templatesFolder = fmt.Sprintf("%s/%s", strings.TrimRight(testsFolder, "/"), "templates")
	}

	return &FSloader{
		TestsFolder:     testsFolder,
		TemplatesFolder: templatesFolder,
	}, nil
}

func (ldr *FSloader) ListTestCases() ([]string, error) {

	_, err := os.Stat(ldr.TestsFolder)
	if err != nil {
		return nil, err
	}

	// No need to check for errors here cause of Glob implementation
	files, _ := filepath.Glob(fmt.Sprintf("%s/*.yaml", ldr.TestsFolder))

	return files, nil
}

func (ldr *FSloader) LoadTestCase(fname string) (*TestCaseSpec, error) {

	var testcase *TestCaseSpec

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

func (ldr *FSloader) LoadTemplate(tname string) (*TemplateSpec, error) {

	var testTempl *TemplateSpec

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
