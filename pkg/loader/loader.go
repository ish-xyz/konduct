package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
)

func NewLoader(ltype, src string) Loader {
	return &FSloader{
		Source: src,
		logger: logrus.New().WithField("name", "fsloader"),
	}
}

func (ldr *FSloader) GetCases() ([]string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.yaml", ldr.Source))
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (ldr *FSloader) Load(files []string) ([]*TestCase, error) {

	tests := []*TestCase{}

	for _, fname := range files {
		var testcase *TestCase
		ldr.logger.Infof("loading file %s", fname)
		data, err := os.ReadFile(fname)
		if err != nil {
			ldr.logger.Warningf("error while reading file %s, error: %v", fname, err)
			continue
		}

		err = yaml.Unmarshal(data, &testcase)
		if err != nil {
			ldr.logger.Warningf("error while loading file %s, error: %v", fname, err)
			continue
		}
		tests = append(tests, testcase)
		ldr.logger.Infof("loaded testcase %s", testcase.Name)
	}

	return tests, nil
}
