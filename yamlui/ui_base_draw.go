package yamlui

// **********************************************************************
// 描画周り
// ==================================================
// 描画インターフェース
type DrawIF interface {
	Draw(x, y float64, ctx DrawContext)
}

// 直接DrawIFを呼び出すのではなく、DrawTreeの中でDrawQueueItemにしてキューに入れる
type DrawQueueItem struct {
	drawIF DrawIF
	x      float64
	y      float64
	z      int
	ctx    DrawContext
}

// Align操作が必要な時とか、DrawTreeを自前で実装したいときのインターフェース
type DrawTreeIF interface {
	DrawTree(lib *YAMLUI, x, y float64, z int, ctx DrawContext)
}

type DrawContext struct {
	Parent     *UIBase // 親のUI
	ParentClip Area    // 親の描画領域（クリップされている）
	Base       *UIBase // 基底のUIBaseを入れてXYWH等に直接アクセスできるようにする
	Clip       Area    // 自分が描画できる領域（クリップされている）
}

func (self *UIBase) SetDrawIF(drawIF DrawIF) {
	self.drawIF = drawIF
}

func (self *UIBase) SetDrawTreeIF(drawTreeIF DrawTreeIF) {
	self.drawTreeIF = drawTreeIF
}

// ==================================================
// コンテキスト作成
// 作成時に座標を計算しておく
func (self *UIBase) calcDrawArea(clip Area) Area {
	if self.IsAbs {
		return Area{
			X: self.X,
			Y: self.Y,
			W: self.W,
			H: self.H,
		}.Clip(clip)
	} else {
		return Area{
			X: clip.X + self.X,
			Y: clip.Y + self.Y,
			W: self.W,
			H: self.H,
		}.Clip(clip)
	}
}

func (self *UIBase) calcDrawPos(x, y float64) (float64, float64) {
	if self.IsAbs {
		return self.X, self.Y
	} else {
		return x + self.X, y + self.Y
	}
}

func NewDrawContext(self *UIBase, parent *UIBase, parentClip Area) DrawContext {
	return DrawContext{
		Parent:     parent,
		ParentClip: parentClip,
		Base:       self,
		Clip:       self.calcDrawArea(parentClip),
	}
}

// ==================================================
// DrawTreeの再帰実行
func (self *UIBase) recDrawTree(lib *YAMLUI, x, y float64, z int, ctx DrawContext) {
	// 子供の描画
	for _, child := range self.children {
		if child.IsVisible {
			// 描画座標やクリップ領域等を計算してコンテキストを作る
			childCtx := NewDrawContext(child, self, ctx.Clip)
			childX, childY := child.calcDrawPos(x, y)

			// 描画インターフェースがあれば描画キューに入れる
			if child.drawIF != nil {
				lib.drawQueue = append(lib.drawQueue, DrawQueueItem{
					drawIF: child.drawIF,
					x:      childX,
					y:      childY,
					z:      z,
					ctx:    childCtx,
				})
			}

			// 自前drawTreeがあればそちらを呼び出す
			if child.drawTreeIF != nil {
				child.drawTreeIF.DrawTree(lib, childX, childY, z+1, childCtx)
			} else {
				// なければ再帰的に子供を描画
				child.recDrawTree(lib, childX, childY, z+1, childCtx)
			}
		}
	}
}
