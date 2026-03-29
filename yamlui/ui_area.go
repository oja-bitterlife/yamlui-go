package yamlui

type UIArea struct {
	Base *UIBase
}

func NewUIArea() *UIArea {
	return &UIArea{
		Base: NewUIBase("Area"),
	}
}
