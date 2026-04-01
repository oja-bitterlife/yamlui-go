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

func NewUIWindow(type_ string) *UIWindow {
	window := &UIWindow{
		Base:  NewUIBase(type_),
		State: UIWindowOpening,
	}
	return window
}
