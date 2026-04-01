package yamlui

// ==================================================
// 更新
// ----------------------------------------
// Updateのインターフェース
type OnInitIF interface {
	OnInit(ui *UIBase, ctx UpdateContext) error
}

type UpdateIF interface {
	Update(ui *UIBase, ctx UpdateContext) error
}

type UpdateTreeIF interface {
	UpdateTree(ui *UIBase, ctx UpdateContext) error
}

type UpdateContext struct {
	Parent *UIBase
	frame  int
}

func (self *UIBase) SetOnInitIF(onInitIF OnInitIF) {
	self.onInitIF = onInitIF
}

func (self *UIBase) SetUpdateIF(updateIF UpdateIF) {
	self.updateIF = updateIF
}

func (self *UIBase) SetUpdateTreeIF(updateTreeIF UpdateTreeIF) {
	self.updateTreeIF = updateTreeIF
}

func NewUpdateContext(parent *UIBase, frame int) UpdateContext {
	return UpdateContext{
		Parent: parent,
		frame:  frame,
	}
}

// ----------------------------------------
// UpdateIFの呼び出し
func (self *UIBase) callUpdate(ctx UpdateContext) error {
	var err error

	if self.IsEnable {
		if self.updateIF != nil {
			self.updateIF.Update(self, ctx)
		}

		// UIの更新後スクリプトがあれば走らせる
		if self.script != nil {
			self.storeToVM(self.script.GetVM())
			_, err = self.script.Run()
			self.loadFromVM(self.script.GetVM())
		}
	}

	return err
}

func (self *UIBase) callUpdateTree(ctx UpdateContext) error {
	if self.updateTreeIF != nil {
		return self.updateTreeIF.UpdateTree(self, ctx)
	} else {
		return self.updateTree(ctx)
	}
}

// ----------------------------------------
// Update実行
// 呼び出し口
func (self *UIBase) Update(frame int) error {
	ctx := NewUpdateContext(self, frame)

	var lastErr error
	if self.Frame == 0 {
		// 最初のフレームならOnInitを呼び出す
		if self.onInitIF != nil {
			lastErr = self.onInitIF.OnInit(self, ctx)
		}
	} else {
		// それ以降のフレームならUpdateを呼び出す
		if err := self.callUpdate(ctx); err != nil {
			lastErr = err
		}
	}

	// InitもUpdateもTreeは手繰る
	if err := self.callUpdateTree(ctx); err != nil {
		lastErr = err
	}

	self.Frame++

	return lastErr
}

// 再帰実行
func (self *UIBase) updateTree(ctx UpdateContext) error {
	ctx = NewUpdateContext(self, ctx.frame)

	var lastErr error

	// 先に子のUpdateを全部実行する
	for _, child := range self.children {
		if err := child.callUpdate(ctx); err != nil {
			lastErr = err
		}
	}

	// そのあとTreeを手繰る
	for _, child := range self.children {
		if err := child.callUpdateTree(ctx); err != nil {
			lastErr = err
		}
	}

	return lastErr
}
