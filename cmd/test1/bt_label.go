package main

import "github.com/oja-bitterlife/yamlui-go/yamlui"

type BTLabel struct {
	Base  *yamlui.UILabel
	model *model
}

func NewBTLabel(m *model, text string) *BTLabel {
	label := &BTLabel{
		Base:  yamlui.NewUILabel(text),
		model: m,
	}
	label.Base.Base.SetDrawIF(label)
	return label
}

func (bt *BTLabel) Draw(ui *yamlui.UIBase, clip yamlui.Area) {
	r := rune(ui.X + '0')                                        // 仮の描画: X座標を文字に変換
	bt.model.canvas[clip.Y][clip.X] = Cell{Rune: r, Color: "86"} // 仮の描画
}
