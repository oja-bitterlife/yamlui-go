package main

import (
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type UI struct {
	name string
}

func main() {
	yamlui.NewYamlUI[UI]()
}
