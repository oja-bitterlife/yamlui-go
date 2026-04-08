package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

const TITLE_CLOSING_EVENT = "title_closing"

type BTCloseWin struct {
	uiBlock *yamlui.UIBlock
	lib     *model
}

func NewBTCloseWin(lib *model) *BTCloseWin {
	return &BTCloseWin{
		uiBlock: yamlui.NewUIBlock(),
		lib:     lib,
	}
}

func (self *BTCloseWin) Clone() yamlui.UICloned {
	return &BTCloseWin{
		uiBlock: self.uiBlock.Clone().(*yamlui.UIBlock),
		lib:     self.lib,
	}
}

func (self *BTCloseWin) GetUIBase() *yamlui.UIBase {
	return self.uiBlock.GetUIBase()
}

func (self *BTCloseWin) Setup(type_ string, data script.ValueMap) error {
	if err := self.uiBlock.Setup(type_, data); err != nil { // super call
		return err
	}
	self.GetUIBase().SetDispatchIF(self)
	self.GetUIBase().SetDrawIF(self)
	return nil
}

func (self *BTCloseWin) Dispatch(lib *yamlui.YAMLUI, event string) (string, error) {
	self.uiBlock.StartBlockTimer(lib.Frame, 30) // 30フレームのブロック
	return "", nil
}

func (self *BTCloseWin) Draw(lib *yamlui.YAMLUI, x, y int, clip yamlui.Area) {
	if self.uiBlock.IsStarted() == false {
		return // ブロックしていない
	}

	win := lib.FindByID("win:title")
	if win == nil {
		script.LogErr("BTCloseWin: win:title not found")
	}

	// タイマーが終了したらウィンドウを閉じる
	if self.uiBlock.Timer.IsFinish(lib.Frame) {
		win.Remove = true
	} else {
		win.SetPropNum("close_ratio", self.uiBlock.Timer.Progress(lib.Frame, 100))
	}
}
