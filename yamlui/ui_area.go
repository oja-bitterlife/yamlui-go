package yamlui

import (
	"github.com/oja-bitterlife/yamlui-go/script"
)

type UIArea struct {
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

func NewUIArea(ui *UIBase, data map[string]script.Value) *UIArea {
	area := &UIArea{
		Base: NewUIBase("area"),
	}
	area.Base.drawTreeIF = area

	// デフォルトはMargin=0
	area.Margin = propINum(data, "Margin", 0)
	area.MarginTop = propINum(data, "MarginTop", 0)
	area.MarginBottom = propINum(data, "MarginBottom", 0)
	area.MarginLeft = propINum(data, "MarginLeft", 0)
	area.MarginRight = propINum(data, "MarginRight", 0)
	area.MarginX = propINum(data, "MarginX", 0)
	area.MarginY = propINum(data, "MarginY", 0)

	// Align系のプロパティはデフォルトはfalse
	area.AlignCenter = propBool(data, "AlignCenter", false)
	area.AlignCenterX = propBool(data, "AlignCenterX", false)
	area.AlignCenterY = propBool(data, "AlignCenterY", false)
	area.AlignRight = propBool(data, "AlignRight", false)
	area.AlignBottom = propBool(data, "AlignBottom", false)

	// IsAbsのデフォルトはfalse
	area.IsAbs = propBool(data, "IsAbs", false)

	return area
}

func (self *UIArea) DrawTree(lib *YAMLUI, x, y float64, z int, ctx DrawContext) {
	// 面倒なので先にintにしておく
	parentW, parentH := ctx.ParentClip.WH()
	selfW, selfH := ctx.Base.Area().WH()

	// 絶対座標対応
	if self.IsAbs {
		x = lib.Screen.X
		y = lib.Screen.Y
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

	ctx.Base.RecDrawTree(lib, alignX+left, alignY+top, z, ctx)
}
