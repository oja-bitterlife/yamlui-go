package yamlui

// ウインドウの状態
type UIWindowState int

const (
	UIWindowOpening UIWindowState = iota
	UIWindowOpened
	UIWindowClosing
	UIWindowClosed
)

type UIWindow struct {
	Base  *UIBase
	State UIWindowState
}

func NewUIWindow() *UIWindow {
	window := &UIWindow{
		Base:  NewUIBase("UIWindow"),
		State: UIWindowOpening,
	}
	return window
}
