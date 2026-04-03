package yamlui

import "slices"

// **********************************************************************
// 描画インターフェース
type DrawContext struct {
	Lib        *YAMLUI // ライブラリ全体へのアクセス
	Clip       Area    // 自分が描画できる領域（クリップされている）
	ParentClip Area    // 親の描画領域（クリップされている）
}

// 直接DrawIFを呼び出すのではなく、DrawTreeの中でDrawQueueItemにしてキューに入れる
type DrawQueueItem struct {
	drawIF DrawIF
	x      float64
	y      float64
	z      int
	ctx    DrawContext
}

// 通常のDraw
type DrawIF interface {
	Draw(x, y float64, ctx DrawContext)
}

// Align操作が必要な時とか、DrawTreeを自前で実装したいときのインターフェース
type DrawTreeIF interface {
	DrawTree(z int, x, y float64, ctx DrawContext)
}

func (self *UIBase) SetDrawIF(drawIF DrawIF) {
	self.drawIF = drawIF
}

func (self *UIBase) SetDrawTreeIF(drawTreeIF DrawTreeIF) {
	self.drawTreeIF = drawTreeIF
}

// コンテキスト作成
func NewDrawContext(lib *YAMLUI, self *UIBase, parentClip Area) DrawContext {
	return DrawContext{
		Lib:        lib,
		ParentClip: parentClip,
		Clip: Area{
			X: parentClip.X + self.X,
			Y: parentClip.Y + self.Y,
			W: self.W,
			H: self.H,
		}.Clip(parentClip),
	}
}

// **********************************************************************
// 呼び出し口
func (self *YAMLUI) Draw(screen Area) {
	self.Screen = screen

	// drawQueueをクリア
	self.drawQueue = []DrawQueueItem{}

	// 描画コンテキストを作成してDrawTreeを呼び出す
	// ----------------------------------------
	ctx := NewDrawContext(self, self.root, screen)
	self.root.RecDrawTree(0, 0, 0, ctx)

	// drawQueueに溜まった描画命令を実行する
	// ----------------------------------------
	// Z順でソートする
	slices.SortStableFunc(self.drawQueue, func(a, b DrawQueueItem) int {
		return a.z - b.z
	})
	// ソートされたqueueを順番に実行する
	for _, item := range self.drawQueue {
		item.drawIF.Draw(item.x, item.y, item.ctx)
	}
}

// **********************************************************************
// DrawTreeの再帰実行
func (self *UIBase) RecDrawTree(z int, x, y float64, ctx DrawContext) {
	// 子供の描画
	for _, child := range self.children {
		if child.IsVisible {
			// 描画座標やクリップ領域等を計算してコンテキストを作る
			childCtx := NewDrawContext(ctx.Lib, child, ctx.Clip)
			childX, childY := x+child.X, y+child.Y

			// 描画インターフェースがあれば描画キューに入れる
			if child.drawIF != nil {
				ctx.Lib.drawQueue = append(ctx.Lib.drawQueue, DrawQueueItem{
					drawIF: child.drawIF,
					x:      childX,
					y:      childY,
					z:      z,
					ctx:    childCtx,
				})
			}

			// 自前drawTreeがあればそちらを呼び出す
			if child.drawTreeIF != nil {
				child.drawTreeIF.DrawTree(z+1, childX, childY, childCtx)
			} else {
				// なければ再帰的に子供を描画
				child.RecDrawTree(z+1, childX, childY, childCtx)
			}
		}
	}
}
