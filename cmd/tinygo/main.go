package main

import (
	_ "embed"

	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

//go:embed ui.json
var uiData []byte // 3. 変数を定義（stringでもOK）

func main() {
	lib := yamlui.NewYAMLUI()
	if err := lib.Load(uiData); err != nil {
		panic(err)
	}

	vm, err := lib.GetRoot().GetScriptVM()
	if err != nil {
		panic(err)
	}
	vm.Run()

	print("UI loaded successfully!")
}
