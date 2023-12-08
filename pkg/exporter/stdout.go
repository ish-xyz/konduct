package exporter

import (
	"fmt"
	"strings"
)

func NewStdoutExporter(printall bool) Exporter {
	return &StdoutExporter{
		PrintAll: printall,
	}
}

func (e *StdoutExporter) Export(r *Report) error {

	fmt.Printf("\n\n\n")
	if r.Status {
		fmt.Printf("tests succeded!\ncompleted tests: %d/%d\n", r.Succeded, r.Failed+r.Succeded)
	} else {
		fmt.Printf("tests failed!\ncompleted tests: %d/%d\n", r.Succeded, r.Failed+r.Succeded)
	}

	fmt.Println("results:")

	for _, res := range r.Results {
		fmt.Println(">>", res.FilePath)
		if !res.Status || e.PrintAll {
			fmt.Printf("	%s\n", strings.TrimSpace(res.Message))
		}
		fmt.Println("--")
	}

	return nil
}

func (e *StdoutExporter) IsVerbose() bool {
	return true
}
