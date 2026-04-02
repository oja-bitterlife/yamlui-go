package yamlui

import "github.com/oja-bitterlife/yamlui-go/script"

// 何もしないUI
// 仮構築用

type UIEmpty struct {
	UIBase *UIBase
}

func NewUIEmpty() *UIEmpty {
	return &UIEmpty{
		UIBase: NewUIBase(),
	}
}

// ==================================================
// UIComponentIFの実装
func (self *UIEmpty) GetUIBase() *UIBase {
	return self.UIBase
}

func (self *UIEmpty) Clone() UIComponent[*UIBase] {
	return &UIEmpty{
		UIBase: self.UIBase.Clone(),
	}
}

func (self *UIEmpty) Setup(lib *YAMLUI, type_ string, parent *UIBase, data map[string]script.Value) error {
	return self.UIBase.Setup(lib, type_, parent, data) // super call
}
