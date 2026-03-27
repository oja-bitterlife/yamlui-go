package main

import (
	"fmt"

	"github.com/oja-bitterlife/yamlui-go/script"
)

func main() {
	// code := "(_.set @y (_.add @y 10.5))"
	src := `(set @msg 'Hello "Ozaki"\'s World')`

	p := script.NewParser(src)
	v, err := p.Parse()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("%+v\n", v)
}
