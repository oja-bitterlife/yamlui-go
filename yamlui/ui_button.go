package yamlui

type UIButton struct {
	Base *UIBase
}

func NewUIButton(onAction func(action string, button *UIButton)) *UIButton {
	button := &UIButton{
		Base: NewUIBase("Button"),
	}

	// Actionがセットされたときの処理
	button.Base.updateFunc = func(base *UIBase, frame int) {
		if base.Action != "" {
			onAction(base.Action, nil)
		}
	}

	return button
}
