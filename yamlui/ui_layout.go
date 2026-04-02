package yamlui

import (
	"github.com/oja-bitterlife/yamlui-go/script"
)

type UILayout struct {
	UIBase *UIBase

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

func NewUILayout() *UILayout {
	return &UILayout{
		UIBase: NewUIBase(),
	}
}

// **********************************************************************
// UIComponentIFの実装
func (self *UILayout) GetUIBase() *UIBase {
	return self.UIBase
}

func (self *UILayout) Clone() *UILayout {
	newLayout := *self
	newLayout.UIBase = self.UIBase.Clone()
	return &newLayout
}

func (self *UILayout) Setup(lib *YAMLUI, type_ string, parent *UIBase, data map[string]script.Value) error {
	self.UIBase.SetDrawTreeIF(self)

	// デフォルトはMargin=0
	self.Margin = PropINum(data, "Margin", 0)
	self.MarginTop = PropINum(data, "MarginTop", 0)
	self.MarginBottom = PropINum(data, "MarginBottom", 0)
	self.MarginLeft = PropINum(data, "MarginLeft", 0)
	self.MarginRight = PropINum(data, "MarginRight", 0)
	self.MarginX = PropINum(data, "MarginX", 0)
	self.MarginY = PropINum(data, "MarginY", 0)

	// Align系のプロパティはデフォルトはfalse
	self.AlignCenter = PropBool(data, "AlignCenter", false)
	self.AlignCenterX = PropBool(data, "AlignCenterX", false)
	self.AlignCenterY = PropBool(data, "AlignCenterY", false)
	self.AlignRight = PropBool(data, "AlignRight", false)
	self.AlignBottom = PropBool(data, "AlignBottom", false)

	// IsAbsのデフォルトはfalse
	self.IsAbs = PropBool(data, "IsAbs", false)

	return nil
}

// **********************************************************************
// DrawTreeのOverride
func (self *UILayout) DrawTree(z int, x, y float64, ctx DrawContext) {
	// 面倒なので先にintにしておく
	parentW, parentH := ctx.ParentClip.WH()
	selfW, selfH := ctx.Base.Area().WH()

	// 絶対座標対応
	if self.IsAbs {
		x = ctx.Lib.Screen.X
		y = ctx.Lib.Screen.Y
	}

	// Align系のプロパティを考慮して、子の描画領域を計算する
	offsetX := 0.0
	offsetY := 0.0
	if self.AlignRight {
		offsetX = parentW - selfW
	}
	if self.AlignBottom {
		offsetY = parentH - selfH
	}
	if self.AlignCenter || self.AlignCenterX {
		offsetX = (parentW - selfW) / 2
	}
	if self.AlignCenter || self.AlignCenterY {
		offsetY = (parentH - selfH) / 2
	}
	// alignされた座標
	alignX := x + offsetX
	alignY := y + offsetY

	// マージンを考慮して、子の描画領域を計算する
	left := float64(self.MarginLeft + self.MarginX + self.Margin)
	top := float64(self.MarginTop + self.MarginY + self.Margin)
	right := float64(self.MarginRight + self.MarginX + self.Margin)
	bottom := float64(self.MarginBottom + self.MarginY + self.Margin)

	// Right/Bottomマージン分親のエリアから引いておく
	parentArea := ctx.ParentClip
	parentArea.W -= right
	parentArea.H -= bottom

	// Left/Topマージン分座標をずらす
	ctx.Clip = Area{
		X: alignX + left,
		Y: alignY + top,
		W: selfW,
		H: selfH,
	}.Clip(parentArea)

	// 元のDrawTreeを呼び出す
	ctx.Base.RecDrawTree(z, alignX+left, alignY+top, ctx)
}
