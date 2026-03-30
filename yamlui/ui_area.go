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

func NewUIArea() *UIArea {
	area := &UIArea{
		Base: NewUIBase("Area"),
	}
	area.Base.drawTreeIF = area
	return area
}

func (self *UIArea) DrawTree(ui *UIBase, clip Area) {
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
