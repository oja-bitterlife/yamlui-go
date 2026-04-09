package yamlui

type UIBlock struct {
	UIBase    *UIBase
	Timer     Timer
	isStarted bool
}

func NewUIBlock(lib *YAMLUI) *UIBlock {
	return &UIBlock{
		UIBase: NewUIBase(lib),
	}
}

// ==================================================
// setter/getter
// イベントを指定してブロックする。引数なしならすべてのイベントをブロックする
func (self *UIBlock) StartBlock(args ...string) {
	self.isStarted = true

	if len(args) == 0 {
		args = []string{"*"}
	}
	self.GetUIBase().Events = self.GetUIBase().Events[:0] // いったん空にする
	self.GetUIBase().Events = append(self.GetUIBase().Events, args...)
}

func (self *UIBlock) StartBlockTimer(startFrame, duration int, args ...string) {
	self.Timer = NewTimer(startFrame, duration)
	self.StartBlock(args...)
}

func (self *UIBlock) EndBlock() {
	self.isStarted = false
	self.GetUIBase().Events = self.GetUIBase().Events[:0] // ブロックしているイベントをクリア
}

func (self *UIBlock) IsStarted() bool {
	return self.isStarted
}

// **********************************************************************
// UIComponentIFの実装
func (self *UIBlock) GetUIBase() *UIBase {
	return self.UIBase
}

func (self *UIBlock) Clone() UICloned {
	return &UIBlock{
		UIBase: self.GetUIBase().Clone(),
	}
}

func (self *UIBlock) Setup(type_ string) {}
