package yamlui

// 何もしないUI
// 仮構築用

type UIEmpty struct {
	UIBase *UIBase
}

func NewUIEmpty(lib *YAMLUI) *UIEmpty {
	return &UIEmpty{
		UIBase: NewUIBase(lib),
	}
}

// ==================================================
// UIComponentIFの実装
func (self *UIEmpty) GetUIBase() *UIBase {
	return self.UIBase
}

func (self *UIEmpty) Clone() UICloned {
	return &UIEmpty{
		UIBase: self.UIBase.Clone(),
	}
}

func (self *UIEmpty) Setup(type_ string) {}
