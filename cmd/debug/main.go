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
	builtin.SetBuiltinCmds(vm)
	result, err := vm.Run(src)
	if err != nil {
		fmt.Printf("{\"Error\":\"%s\"}\n", err.Error())
		return
	}
	result.Dump()

	// vmのcmsを表示
	fmt.Println("Commands in VM:")
	for cmd := range vm.GetCmds() {
		fmt.Printf("%s\n", cmd)
	}
}
