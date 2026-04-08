package main

import (
	_ "embed"

	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

//go:embed ui.json
var uiData []byte // 3. 変数を定義（stringでもOK）

func main() {
	lib := yamlui.NewYAMLUI()
	// if err := lib.Load(uiData); err != nil {
	// 	panic(err)
	// }

	// vm := lib.GetRoot().GetVM()
	// vm.Run()

	if err := lib.Start(uiData); err != nil {
		panic(err)
	}
	lib.Draw(0, 0, 800, 600)
	lib.Stop()

	print("UI loaded successfully!")
}
