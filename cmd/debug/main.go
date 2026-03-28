package main

import (
	"fmt"

	"github.com/oja-bitterlife/yamlui-go/script"
)

func main() {
	// code := "(_.set @y (_.add @y 10.5))"
	// src := `set @msg 'Hello "Ozaki"\'s World'`
	src := `
(set _count 3)
(repeat _i _count
  (switch _i "First!" "Second!" "Third!"))
`

	vm := script.NewVM()
	result, err := vm.Run(src)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	result.Dump()
}
