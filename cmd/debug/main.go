package main

import (
	"fmt"

	"github.com/oja-bitterlife/yamlui-go/script"
)

func main() {
	code := "(_.set @y (_.add @y 10.5))"
	ast := script.Parse(code)

	// コンパイル
	insts := script.Compile(ast)

	fmt.Println("--- Compiled Instructions ---")
	for i, ins := range insts {
		// OpCode をとりあえず数値で表示
		fmt.Printf("%03d: Op=%d, Val.Str=%s, Val.Num=%g\n", i, ins.Op, ins.Val.Str, ins.Val.Number)
	}
}
