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

func (self *UIWindow) Setup(lib *YAMLUI, type_ string, parent *UIBase, data map[string]script.Value) error {
	return self.UIBase.Setup(lib, type_, parent, data) // super call
}
