package main

import (
	"github.com/ish-xyz/ykubetest/cmd"
)

//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen crd object:headerFile="hack/boilerplate.go.txt" paths="./crdgen/..."
//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen crd:crdVersions=v1,allowDangerousTypes=true paths=./crdgen/... output:dir="charts/kubetest-crd"
//go:generate rm -rf ./config
func main() {
	cmd.Execute()
}
