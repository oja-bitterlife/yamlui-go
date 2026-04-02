package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTLabel struct {
	UIBase *yamlui.UIBase
	model  *model
}

func NewBTLabel() *BTLabel {
	return &BTLabel{
		UIBase: yamlui.NewUIBase(),
	}
}

func (self *BTLabel) GetUIBase() *yamlui.UIBase {
	return self.UIBase
}

func (self *BTLabel) Clone() yamlui.UIComponent[*yamlui.UIBase] {
	return &BTLabel{
		UIBase: self.UIBase.Clone(),
		model:  self.model,
	}
}

func (self *BTLabel) Setup(lib *yamlui.YAMLUI, type_ string, parent *yamlui.UIBase, data map[string]script.Value) error {
	self.UIBase.SetDrawIF(self)
	return nil
}

func (self *BTLabel) Draw(x, y float64, ctx yamlui.DrawContext) {
	panic("BTLabel.Draw is not implemented yet")
}
