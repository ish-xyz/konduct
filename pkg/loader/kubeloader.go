package loader

import (
	"context"
	"fmt"

	"github.com/ish-xyz/ykubetest/pkg/client"
	"github.com/mitchellh/mapstructure"
)

func NewKubeLoader(cl client.Client) (Loader, error) {
	return &KubeLoader{
		Client: cl,
	}, nil
}

func (ldr *KubeLoader) ListTestCases() ([]string, error) {
	var tcases []string

	resp := ldr.Client.Get(context.TODO(), "kubetest.io/v1", "testcase", "", "", "")
	if resp.Error != "" {
		return nil, fmt.Errorf("%v", resp.Error)
	}
	for _, c := range resp.Objects {
		tcases = append(tcases, c["metadata"].(map[string]interface{})["name"].(string))
	}

	return tcases, nil
}

// TODO join the following 2 functions with generics

func (ldr *KubeLoader) LoadTestCase(resourceName string) (*TestCase, error) {

	var testcase *TestCase

	resp := ldr.Client.Get(context.TODO(), "kubetest.io/v1", "testcase", "", resourceName, "")
	if resp.Error != "" {
		return nil, fmt.Errorf("%v", resp.Error)
	}

	if len(resp.Objects) == 0 {
		return nil, fmt.Errorf("testcase not found")
	}

	err := mapstructure.Decode(resp.Objects[0]["spec"], &testcase)

	return testcase, err
}

func (ldr *KubeLoader) LoadTemplate(resourceName string) (*Template, error) {
	var templ *Template

	resp := ldr.Client.Get(context.TODO(), "kubetest.io/v1", "template", "", resourceName, "")
	if resp.Error != "" {
		return nil, fmt.Errorf("%v", resp.Error)
	}

	if len(resp.Objects) == 0 {
		return nil, fmt.Errorf("template not found")
	}

	err := mapstructure.Decode(resp.Objects[0]["spec"], &templ)

	return templ, err
}
