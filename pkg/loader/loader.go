package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

func NewLoader(ltype, src string) Loader {
	return &FSloader{
		Source: src,
	}
}

func (ldr *FSloader) GetTestCases() ([]string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.yaml", ldr.Source))
	if err != nil {
		return nil, err
	}

	return files, nil
}

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
