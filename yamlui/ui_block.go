package yamlui

import (
	"github.com/oja-bitterlife/yamlui-go/script"
)

type UIBlock struct {
	UIBase *UIBase
}

func NewUIBlock() *UIBlock {
	return &UIBlock{
		UIBase: NewUIBase(),
	}
}

// ==================================================
// setter/getter
// イベントを指定してブロックする。引数なしならすべてのイベントをブロックする
func (self *UIBlock) SetBlock(args ...string) {
	if len(args) == 0 {
		args = []string{"*"}
	}
	self.GetUIBase().Events = args
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
	return nil
}
