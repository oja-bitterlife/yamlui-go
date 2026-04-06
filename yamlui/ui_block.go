package yamlui

import (
	"github.com/oja-bitterlife/yamlui-go/script"
)

type UIBlock struct {
	UIBase     *UIBase
	timer      Timer
	timerEvent string
	actionFunc func(lib *YAMLUI, event string) error
}

func NewUIBlock() *UIBlock {
	return &UIBlock{
		UIBase: NewUIBase(),
	}
}

// ==================================================
// setter
func (self *UIBlock) SetActionFunc(actionFunc func(lib *YAMLUI, event string) error) {
	self.actionFunc = actionFunc
}

func (self *UIBlock) SetTimer(event string, duration int) {
	self.timer = NewTimer(0, duration)
	self.timerEvent = event
}

// **********************************************************************
// UIComponentIFの実装
func (self *UIBlock) GetUIBase() *UIBase {
	return self.UIBase
}

func (self *UIBlock) Clone() UICloned {
	return &UIBlock{
		UIBase: self.GetUIBase().Clone(),
	}
}

func (self *UIBlock) Setup(type_ string, data script.ValueMap) error {
	if err := self.GetUIBase().Setup(type_, data); err != nil { // super call
		return err
	}

	// すべてのイベントをこのUIBlockで吸収する
	self.GetUIBase().Events = []string{"*"}

	// イベントチェック用Update
	self.GetUIBase().SetUpdateIF(self)
	return nil
}

// **********************************************************************
// Updateでイベントを吸収する
func (self *UIBlock) Update(lib *YAMLUI, events []string) error {
	// タイマーイベントチェック
	if self.timer.Duration() > 0 && self.timerEvent != "" {
		if self.timer.IsFinish(self.GetUIBase().UpdateCount) {
			if self.actionFunc != nil {
				self.actionFunc(lib, self.timerEvent)
			}
		}
	}

	// それ以外のイベントも全部回す.PressAnyKey的な使い方を想定
	for _, event := range events {
		self.actionFunc(lib, event)
	}

	return nil
}
