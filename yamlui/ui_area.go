package yamlui

import "github.com/oja-bitterlife/yamlui-go/script"

type UIArea struct {
	Base         *UIBase
	Margin       int
	MarginTop    int
	MarginBottom int
	MarginLeft   int
	MarginRight  int
	MarginX      int
	MarginY      int
}

func NewUIArea(type_ string, ui *UIBase, data map[string]script.Value) *UIArea {
	area := &UIArea{
		Base: NewUIBase(type_),
	}
	area.Base.drawTreeIF = area

	area.Margin = propINum(data, "Margin", 0)
	area.MarginTop = propINum(data, "MarginTop", 0)
	area.MarginBottom = propINum(data, "MarginBottom", 0)
	area.MarginLeft = propINum(data, "MarginLeft", 0)
	area.MarginRight = propINum(data, "MarginRight", 0)
	area.MarginX = propINum(data, "MarginX", 0)
	area.MarginY = propINum(data, "MarginY", 0)

	return area
}

func (self *UIArea) DrawTree(ui *UIBase, clip Area) {

	// マージンを考慮して、子の描画領域を計算する
	left := self.MarginLeft + self.MarginX + self.Margin
	top := self.MarginTop + self.MarginY + self.Margin
	right := self.MarginRight + self.MarginX + self.Margin
	bottom := self.MarginBottom + self.MarginY + self.Margin

	area := Area{
		X: clip.X + float64(left),
		Y: clip.Y + float64(top),
		W: clip.W - float64(left+right),
		H: clip.H - float64(top+bottom),
	}

	ui.drawTree(area)
}
