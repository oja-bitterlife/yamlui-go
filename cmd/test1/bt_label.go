package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTLabel struct {
	Base  *yamlui.UILabel
	model *model
}

func NewBTLabel(m *model, parent *yamlui.UIBase, data map[string]script.Value) *BTLabel {
	label := &BTLabel{
		Base:  yamlui.NewUILabel("test"),
		model: m,
	}
	label.Base.Base.SetDrawIF(label)
	return label
}

func (self *BTLabel) Draw(x, y float64, ctx yamlui.DrawContext) {
	drawY := int(y)
	if ctx.Clip.IContainsY(drawY) == false {
		return
	}

	for j, char := range ctx.Base.Text {
		drawX := int(x) + j
		if ctx.Clip.IContainsX(drawX) == false {
			continue
		}
		self.model.canvas[drawY][drawX] = Cell{Rune: char, Color: "white"}
	}

}
