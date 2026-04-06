package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTWinClose struct {
	uiBlock *yamlui.UIBlock
	lib     *model
}

func NewBTWinClose(lib *model) *BTWinClose {
	return &BTWinClose{
		uiBlock: yamlui.NewUIBlock(),
		lib:     lib,
	}
}

func (self *BTWinClose) Clone() yamlui.UICloned {
	return &BTWinClose{
		uiBlock: self.uiBlock.Clone().(*yamlui.UIBlock),
		lib:     self.lib,
	}
}

func (self *BTWinClose) GetUIBase() *yamlui.UIBase {
	return self.uiBlock.GetUIBase()
}

func (self *BTWinClose) Setup(type_ string, data script.ValueMap) error {
	if err := self.uiBlock.Setup(type_, data); err != nil { // super call
		return err
	}
	return nil
}
