package yamlui

type UIButton struct {
	Base *UIBase
}

func (button *UIButton) Update(ui *UIBase, frame int) {
	// Actionがセットされたときの処理
	if ui.ScriptReturn != "" {
	}
}

func NewUIButton(onAction func(action string, button *UIButton)) *UIButton {
	button := &UIButton{
		Base: NewUIBase("Button"),
	}

	return button
}
