package loader

import "github.com/ish-xyz/ykubetest/pkg/client"

func NewKubeLoader(cl client.Client) (Loader, error) {
	return &KubeLoader{
		Client: cl,
	}, nil
}

func (ldr *KubeLoader) ListTestCases() ([]string, error) {
	var files []string
	return files, nil
}

// TODO join the following 2 functions with generics

func (ldr *KubeLoader) LoadTestCase(resourceName string) (*TestCase, error) {

	var testcase *TestCase

	return testcase, nil
}

func (ldr *KubeLoader) LoadTemplate(tname string) (*TestTemplate, error) {

	var testTempl *TestTemplate

	return testTempl, nil
}
