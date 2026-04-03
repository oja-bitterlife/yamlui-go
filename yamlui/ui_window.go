package yamlui

import "github.com/oja-bitterlife/yamlui-go/script"

// ウインドウの状態
type UIWindowState int

const (
	UIWindowOpening UIWindowState = iota
	UIWindowOpened
	UIWindowClosing
	UIWindowClosed
)

type UIWindow struct {
	UIBase *UIBase
	State  UIWindowState
}

func NewUIWindow() *UIWindow {
	window := &UIWindow{
		UIBase: NewUIBase(),
		State:  UIWindowOpening,
	}
	return window
}

// **********************************************************************
// UIComponentIFの実装
func (self *UIWindow) GetUIBase() *UIBase {
	return self.UIBase
}

func (self *UIWindow) Clone() UIComponent[*UIBase] {
	return &UIWindow{
		UIBase: self.UIBase.Clone(),
		State:  self.State,
	}
}

func (self *UIWindow) Setup(type_ string, data script.ValueMap) error {
	return self.UIBase.Setup(type_, data) // super call
}
