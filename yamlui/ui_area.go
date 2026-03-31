package yamlui

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

func NewUIArea(ui *UIBase, data map[string]any) *UIArea {
	area := &UIArea{
		Base: NewUIBase("Area"),
	}
	area.Margin = propNum(data, "Margin", 0)
	area.MarginTop = propNum(data, "MarginTop", 0)
	area.MarginBottom = propNum(data, "MarginBottom", 0)
	area.MarginLeft = propNum(data, "MarginLeft", 0)
	area.MarginRight = propNum(data, "MarginRight", 0)
	area.MarginX = propNum(data, "MarginX", 0)
	area.MarginY = propNum(data, "MarginY", 0)
	area.Base.drawTreeIF = area
	return area
}

func (self *UIArea) DrawTree(ui *UIBase, clip Area) {
	// マージンを考慮して、子の描画領域を計算する
	left := self.MarginLeft + self.MarginX + self.Margin
	top := self.MarginTop + self.MarginY + self.Margin
	right := self.MarginRight + self.MarginX + self.Margin
	bottom := self.MarginBottom + self.MarginY + self.Margin

	area := Area{
		X: clip.X + left,
		Y: clip.Y + top,
		W: clip.W - left - right,
		H: clip.H - top - bottom,
	}

	ui.drawTree(area)
}
