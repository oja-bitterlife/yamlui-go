package yamlui

// **********************************************************************
// 描画周り
// ==================================================
// 描画インターフェース
type DrawIF interface {
	Draw(x, y float64, ctx DrawContext)
}

// Align操作が必要な時とか、DrawTreeを自前で実装したいときのインターフェース
type DrawTreeIF interface {
	DrawTree(x, y float64, ctx DrawContext)
}

type DrawContext struct {
	Parent     *UIBase // 親のUI
	ParentClip Area    // 親の描画領域（クリップされている）
	State      *UIBase // 基底のUIBaseを入れてXYWH等に直接アクセスできるようにする
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
		State:      self,
		Clip:       self.calcDrawArea(parentClip),
	}
}

// ----------------------------------------
// 描画IFの呼び出し
func (self *UIBase) callDraw(x, y float64, ctx DrawContext) {
	if self.IsVisible && self.drawIF != nil {
		x, y := self.calcDrawPos(ctx.Parent.X, ctx.Parent.Y)
		self.drawIF.Draw(x, y, ctx)
	}
}

func (self *UIBase) callDrawTree(x, y float64, ctx DrawContext) {
	if self.drawTreeIF != nil {
		self.drawTreeIF.DrawTree(x, y, ctx)
	} else {
		self.drawTree(x, y, ctx)
	}
}

// ----------------------------------------

// 呼び出し口
func (self *UIBase) Draw(screen Area) {
	ctx := NewDrawContext(self, nil, screen)

	// 先に自分を描画する
	self.callDraw(0, 0, ctx)

	// 子のDrawを呼び出す
	self.callDrawTree(0, 0, ctx)
}

// 再帰実行
func (self *UIBase) drawTree(x, y float64, ctx DrawContext) {

	// 先に子のDrawを全部実行する
	for _, child := range self.children {
		if child.drawIF != nil && child.IsVisible {
			childCtx := NewDrawContext(child, self, ctx.Clip)
			childX, childY := child.calcDrawPos(x, y)
			child.callDraw(childX, childY, childCtx)
		}
	}

	// そのあとTreeを手繰る
	for _, child := range self.children {
		if child.IsVisible {
			childCtx := NewDrawContext(child, self, ctx.Clip)
			childX, childY := child.calcDrawPos(x, y)
			child.callDrawTree(childX, childY, childCtx)
		}
	}
}
