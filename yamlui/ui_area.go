package yamlui

type UIArea struct {
	Base         *UIBase
	margin       int
	marginTop    int
	marginBottom int
	marginLeft   int
	marginRight  int
	marginX      int
	marginY      int
}

func NewUIArea() *UIArea {
	area := &UIArea{
		Base: NewUIBase("Area"),
	}
	area.Base.drawTreeIF = area
	return area
}

func (self *UIArea) DrawTree(ui *UIBase, clip Area) {
	left := self.marginLeft + self.marginX + self.margin
	top := self.marginTop + self.marginY + self.margin
	right := self.marginRight + self.marginX + self.margin
	bottom := self.marginBottom + self.marginY + self.margin

	area := Area{
		X: clip.X + left,
		Y: clip.Y + top,
		W: clip.W - left - right,
		H: clip.H - top - bottom,
	}

	ui.drawTree(area)
}
