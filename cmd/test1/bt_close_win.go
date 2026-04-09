package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

const TITLE_CLOSING_EVENT = "title_closing"

type BTCloseWin struct {
	uiBlock *yamlui.UIBlock
	model   *model
}

func NewBTCloseWin(lib *yamlui.YAMLUI, model *model) *BTCloseWin {
	return &BTCloseWin{
		uiBlock: yamlui.NewUIBlock(lib),
		model:   model,
	}
}

func (self *BTCloseWin) Clone() yamlui.UICloned {
	return &BTCloseWin{
		uiBlock: self.uiBlock.Clone().(*yamlui.UIBlock),
		model:   self.model,
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
	return nil
}

func (self *BTCloseWin) Dispatch(lib *yamlui.YAMLUI, event string) {
	if lib.MatchEvent("next:*", event) {
		self.uiBlock.StartBlockTimer(lib.Frame, 30) // 30フレームのブロック
	}

	if self.uiBlock.IsStarted() == false {
		return // 開始していない場合は何もしない
	} else {

		win := lib.FindByID("win:title")
		if win == nil {
			script.LogErr("BTCloseWin: win:title not found")
		}
		// frame
		script.Log("BTCloseWin: frame=%d", lib.Frame)
		// タイマーが終了したらウィンドウを閉じる
		if self.uiBlock.Timer.IsFinish(lib.Frame) {
			win.Remove = true
		} else {
			win.SetPropNum("close_ratio", self.uiBlock.Timer.Progress(lib.Frame, 100))
		}
	}

	// bubble tea に更新を通知して再描画させる
	self.model.NextCmd = self.model.BubbleTeaUpdate()
	script.LogErr("dispatch: %p", self.model)
}
