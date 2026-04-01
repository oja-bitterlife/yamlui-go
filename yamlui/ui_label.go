package yamlui

type UILabel struct {
	Base *UIBase
}

func NewUILabel(text string) *UILabel {
	label := &UILabel{
		Base: NewUIBase("UILabel"),
	}
	label.Base.Text = text
	return label
}
