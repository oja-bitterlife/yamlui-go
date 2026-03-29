package main

import "github.com/oja-bitterlife/yamlui-go/yamlui"

type BTLabel struct {
	Base   *yamlui.UILabel
	canvas [][]Cell
}

func NewBTLabel(text string) *BTLabel {
	label := &BTLabel{
		Base: yamlui.NewUILabel(text),
	}
	label.Base.Base.SetDrawIF(label)
	return label
}

func (bt *BTLabel) Draw(ui *yamlui.UIBase, x, y int) {
	ui.X = x
	ui.Y = y
	bt.canvas[y][x] = Cell{Rune: 'L', Color: "white"} // 仮の描画
}
