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

	return area
}

func (self *UIArea) DrawTree(ui *UIBase, clip Area, ctx DrawContext) {
	// 面倒なので先にintにしておく
	parentW, parentH := ctx.ParentArea.WH()
	selfW, selfH := self.Base.Area().WH()

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
	alignX := self.Base.X + offsetX
	alignY := self.Base.Y + offsetY

	// マージンを考慮して、子の描画領域を計算する
	left := float64(self.MarginLeft + self.MarginX + self.Margin)
	top := float64(self.MarginTop + self.MarginY + self.Margin)
	right := float64(self.MarginRight + self.MarginX + self.Margin)
	bottom := float64(self.MarginBottom + self.MarginY + self.Margin)

	// Right/Bottomマージン分親のエリアから引いておく
	parentArea := ctx.ParentArea
	parentArea.W -= right
	parentArea.H -= bottom

	// Left/Topマージン分座標をずらす
	area := Area{
		X: alignX + left,
		Y: alignY + top,
		W: selfW,
		H: selfH,
	}.Clip(parentArea)

	ui.drawTree(area)
}
