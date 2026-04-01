package yamlui

// ==================================================
// 更新
// ----------------------------------------
// Updateのインターフェース
type OnInitIF interface {
	OnInit(ctx UpdateContext) error
}

type UpdateIF interface {
	Update(lib *YAMLUI, frame int, z int, ctx UpdateContext) error
}

type UpdateTreeIF interface {
	UpdateTree(lib *YAMLUI, frame int, z int, ctx UpdateContext) error
}

// 直接DrawIFを呼び出すのではなく、UpdateTreeの中でUpdateQueueItemにしてキューに入れる
type UpdateQueueItem struct {
	UpdateIF UpdateIF
	z        int
	ctx      UpdateContext
}

type UpdateContext struct {
	Parent *UIBase
	Base   *UIBase // 基底のUIBaseを入れてXYWH等に直接アクセスできるようにする
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

// ==================================================
// Update実行
// コンテキスト作成
func NewUpdateContext(self *UIBase, parent *UIBase) UpdateContext {
	return UpdateContext{
		Parent: parent,
		Base:   self,
	}
}

// 再帰実行
func (self *UIBase) RecUpdateTree(lib *YAMLUI, frame int, z int, ctx UpdateContext) error {
	var lastErr error

	// 子供の更新
	for _, child := range self.children {
		if child.IsEnable {
			// 更新コンテキストを作成
			childCtx := NewUpdateContext(child, self)

			// 更新インターフェースがあれば更新キューに入れる
			if child.updateIF != nil {
				lib.updateQueue = append(lib.updateQueue, UpdateQueueItem{
					UpdateIF: child.updateIF,
					z:        z,
					ctx:      childCtx,
				})
			}

			// UIの更新後スクリプトがあれば走らせる
			if self.script != nil {
				self.storeToVM(self.script.GetVM())
				if _, err := self.script.Run(); err != nil {
					lastErr = err
				}
				self.loadFromVM(self.script.GetVM())
			}

			// UpdateTreeIFがあればそちらを呼び出す。なければ再帰的に呼び出す
			if self.updateTreeIF != nil {
				if err := self.updateTreeIF.UpdateTree(lib, frame, z+1, ctx); err != nil {
					lastErr = err
				}
			} else {
				if err := self.RecUpdateTree(lib, frame, z+1, ctx); err != nil {
					lastErr = err
				}
			}
		}
	}

	return lastErr
}
