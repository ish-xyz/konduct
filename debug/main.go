package main

import (
	"fmt"

	"github.com/antonmedv/expr"
)

type Payload struct {
	Error error
	Output string
	Objects []interface
}

func main() {

	

	env := map[string]interface{}{
		"payload": map[string]interface{}{
			"something":      1,
			"something-else": 1,
		},
	}

	code := `
		payload
	`

	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}
