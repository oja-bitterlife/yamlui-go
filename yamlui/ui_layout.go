package yamlui

import (
	"github.com/oja-bitterlife/yamlui-go/script"
)

type UILayout struct {
	Base *UIBase
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

func (self *UILayout) GetBase() *UIBase {
	return self.Base
}

func NewUIArea(type_ string, ui *UIBase, data map[string]script.Value) *UILayout {
	area := &UILayout{
		Base: NewUIBase(type_),
	}
	area.Base.drawTreeIF = area

	// デフォルトはMargin=0
	area.Margin = PropINum(data, "Margin", 0)
	area.MarginTop = PropINum(data, "MarginTop", 0)
	area.MarginBottom = PropINum(data, "MarginBottom", 0)
	area.MarginLeft = PropINum(data, "MarginLeft", 0)
	area.MarginRight = PropINum(data, "MarginRight", 0)
	area.MarginX = PropINum(data, "MarginX", 0)
	area.MarginY = PropINum(data, "MarginY", 0)

	// Align系のプロパティはデフォルトはfalse
	area.AlignCenter = PropBool(data, "AlignCenter", false)
	area.AlignCenterX = PropBool(data, "AlignCenterX", false)
	area.AlignCenterY = PropBool(data, "AlignCenterY", false)
	area.AlignRight = PropBool(data, "AlignRight", false)
	area.AlignBottom = PropBool(data, "AlignBottom", false)

	// IsAbsのデフォルトはfalse
	area.IsAbs = PropBool(data, "IsAbs", false)

	return area
}

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

	ctx.Base.RecDrawTree(z, alignX+left, alignY+top, ctx)
}
