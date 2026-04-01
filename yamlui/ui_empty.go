package yamlui

// 何もしないUI
// 仮構築用

type UIEmpty struct {
	Base *UIBase
}

func NewUIEmpty(type_ string) *UIEmpty {
	empty := &UIEmpty{
		Base: NewUIBase(type_),
	}
	return empty
}
