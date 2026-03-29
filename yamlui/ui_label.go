package yamlui

type UILabelBase struct {
	Base *UIBase
}

func NewUILabelBase() *UILabelBase {
	return &UILabelBase{
		Base: NewUIBase("Label"),
	}
}
