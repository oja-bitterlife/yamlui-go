package yamlui

// ==================================================
// 更新
// ----------------------------------------
// Updateのインターフェース
type OnInitIF interface {
	OnInit(ctx UpdateContext) error
}

type UpdateIF interface {
	Update(ctx UpdateContext) error
}

type UpdateTreeIF interface {
	UpdateTree(z int, ctx UpdateContext) error
}

// 直接DrawIFを呼び出すのではなく、UpdateTreeの中でUpdateQueueItemにしてキューに入れる
type UpdateQueueItem struct {
	UpdateIF UpdateIF
	z        int
	ctx      UpdateContext
}

type UpdateContext struct {
	Lib    *YAMLUI
	Parent *UIBase
	Base   *UIBase // 基底のUIBaseを入れてXYWH等に直接アクセスできるようにする
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
func NewUpdateContext(lib *YAMLUI, self *UIBase, parent *UIBase) UpdateContext {
	return UpdateContext{
		Lib:    lib,
		Parent: parent,
		Base:   self,
	}
}

// 再帰実行
func (self *UIBase) RecUpdateTree(z int, ctx UpdateContext) error {
	var lastErr error

	// 子供の更新
	for _, child := range self.children {
		if child.IsEnable {
			// 更新コンテキストを作成
			childCtx := NewUpdateContext(ctx.Lib, child, self)

			// 更新インターフェースがあれば更新キューに入れる
			if child.updateIF != nil {
				ctx.Lib.updateQueue = append(ctx.Lib.updateQueue, UpdateQueueItem{
					UpdateIF: child.updateIF,
					z:        z,
					ctx:      childCtx,
				})
			}

			// UpdateTreeIFがあればそちらを呼び出す。なければ再帰的に呼び出す
			if self.updateTreeIF != nil {
				if err := self.updateTreeIF.UpdateTree(z+1, ctx); err != nil {
					lastErr = err
				}
			} else {
				if err := self.RecUpdateTree(z+1, ctx); err != nil {
					lastErr = err
				}
			}
		}
	}

	return lastErr
}
