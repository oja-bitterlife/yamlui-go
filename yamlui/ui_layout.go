package yamlui

type UILayout struct {
	UIBase *UIBase

	// Layout用のプロパティ.ここに置くとprivate扱いになる
	// publicにしたい場合はGetUIBase().SetPropで保存する
	// ==================================================
	// Margin用
	Margin       int
	MarginTop    int
	MarginBottom int
	MarginLeft   int
	MarginRight  int
	MarginX      int
	MarginY      int

	// Align系のプロパティ
	AlignCenter  bool
	AlignCenterX bool
	AlignCenterY bool
	AlignRight   bool
	AlignBottom  bool

	// 自分のXYは絶対座標
	IsAbs bool
}

func NewUILayout(lib *YAMLUI) *UILayout {
	layout := &UILayout{
		UIBase: NewUIBase(lib),
	}
	return layout
}

// **********************************************************************
// UIComponentIFの実装
func (self *UILayout) GetUIBase() *UIBase {
	return self.UIBase
}

func (self *UILayout) Clone() UICloned {
	newLayout := *self
	newLayout.UIBase = self.UIBase.Clone()
	return &newLayout
}

func (self *UILayout) Dispatch(lib *YAMLUI, event string) {
	if event != EVENT_UI_ONCREATE {
		return
	}

	// デフォルトはMargin=0
	self.Margin = self.UIBase.PropNum("@Margin")
	self.MarginTop = self.UIBase.PropNum("@MarginTop")
	self.MarginBottom = self.UIBase.PropNum("@MarginBottom")
	self.MarginLeft = self.UIBase.PropNum("@MarginLeft")
	self.MarginRight = self.UIBase.PropNum("@MarginRight")
	self.MarginX = self.UIBase.PropNum("@MarginX")
	self.MarginY = self.UIBase.PropNum("@MarginY")

	// Align系のプロパティはデフォルトはfalse
	self.AlignCenter = self.UIBase.PropBool("@AlignCenter")
	self.AlignCenterX = self.UIBase.PropBool("@AlignCenterX")
	self.AlignCenterY = self.UIBase.PropBool("@AlignCenterY")
	self.AlignRight = self.UIBase.PropBool("@AlignRight")
	self.AlignBottom = self.UIBase.PropBool("@AlignBottom")

	// IsAbsのデフォルトはfalse
	self.IsAbs = self.UIBase.PropBool("@IsAbs")
}

// **********************************************************************
// DrawTreeのOverride
func (self *UILayout) DrawTreeFilter(lib *YAMLUI, ctx DrawTreeContext) DrawTreeContext {
	area := self.GetUIBase().Area()
	parentW, parentH := ctx.ParentClip.W, ctx.ParentClip.H

	// 絶対座標対応
	if self.IsAbs {
		ctx.X = lib.Screen.X + area.X
		ctx.Y = lib.Screen.Y + area.Y
	}

	// Align系のプロパティを考慮して、子の描画領域を計算する
	offsetX := 0
	offsetY := 0
	if self.AlignRight {
		offsetX = parentW - area.W
	}
	if self.AlignBottom {
		offsetY = parentH - area.H
	}
	if self.AlignCenter || self.AlignCenterX {
		offsetX = (parentW - area.W) / 2
	}
	if self.AlignCenter || self.AlignCenterY {
		offsetY = (parentH - area.H) / 2
	}
	// alignされた座標
	alignX := ctx.X + offsetX
	alignY := ctx.Y + offsetY

	// マージン
	lMgn := self.MarginLeft + self.MarginX + self.Margin
	tMgn := self.MarginTop + self.MarginY + self.Margin
	rMgn := self.MarginRight + self.MarginX + self.Margin
	bMgn := self.MarginBottom + self.MarginY + self.Margin

	// Right/Bottomマージン分親のエリアから引いておく
	//	parentArea := ctx.ParentClip
	//	parentArea.W -= rMgn
	//	parentArea.H -= bMgn

	// マージン分座標をずらす
	newClip := Area{
		X: alignX + lMgn,
		Y: alignY + tMgn,
		W: area.W - lMgn - rMgn,
		H: area.H - tMgn - bMgn,
	}.Clip(ctx.ParentClip)

	// コンテキストの更新
	ctx.X = alignX + lMgn
	ctx.Y = alignY + tMgn
	ctx.Clip = newClip

	// 更新したコンテキストを返す
	return ctx
}
