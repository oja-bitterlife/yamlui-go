package yamlui

type UIWindow struct {
	Base *UIBase
}

func NewUIWindow() *UIWindow {
	window := &UIWindow{
		Base: NewUIBase("UIWindow"),
	}
	return window
}
