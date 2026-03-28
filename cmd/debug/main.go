package main

import (
	"fmt"

	"github.com/oja-bitterlife/yamlui-go/builtin"
	"github.com/oja-bitterlife/yamlui-go/script"
)

func main() {
	// 	src := `
	// (set _count 3)
	// (repeat _i _count
	// 	(switch _i "First!" "Second!" "Third!"))
	// `
	src := `
(! 1 (+ 2 3) 4)
`
	vm := script.NewVM()
	vm.RegisterCmdList(builtin.MathCmds)
	result, err := vm.Run(src)
	if err != nil {
		fmt.Printf("{\"Error\":\"%s\"}\n", err.Error())
		return
	}
	result.Dump()
}
