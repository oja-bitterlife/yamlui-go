package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTLabel struct {
	Base  *yamlui.UIBase
	model *model
	texts []string
}

func (self *BTLabel) GetBase() *yamlui.UIBase {
	return self.Base
}

func NewBTLabel(componentName string, parent *yamlui.UIBase, data map[string]script.Value) *BTLabel {
	label := &BTLabel{
		Base: yamlui.NewUIBase(componentName),
	}
	label.Base.SetDrawIF(label)
	return label
}

func (self *BTLabel) Draw(x, y float64, ctx yamlui.DrawContext) {
}
