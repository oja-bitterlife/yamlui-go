package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

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
	self.GetUIBase().SetUpdateIF(self)
	return nil
}

func (self *BTCloseWin) Update(lib *yamlui.YAMLUI, events []string) error {
	// "next:*"イベントがあったら30フレーム後にウィンドウを閉じるタイマーをセットする
	if lib.HasEvent("next:*", events) {
		self.uiBlock.StartBlockTimer(30)
		self.GetUIBase().Action = TEA_UPDATE_EVENT
	}

	// タイマーが終わるまでイベントを繋いでおく
	if lib.HasEvent(TEA_UPDATE_EVENT, events) {
		self.GetUIBase().Action = TEA_UPDATE_EVENT
	}

	// タイマーが終了したらウィンドウを閉じる
	if self.uiBlock.IsTimerFinish() {
		win := lib.FindByID("win:title")
		if win != nil {
			win.Remove = true
		}
	}

	return nil
}
