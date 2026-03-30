package main

import (
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

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

func (self *BTLabel) Draw(ui *yamlui.UIBase, clip yamlui.Area) {
	orgArea := clip.Inset(-ui.X, -ui.Y)
	cx, cy := orgArea.AlignCenter(len(ui.Text), 1)
	cx += ui.X
	cy += ui.Y

	for i, r := range ui.Text {
		targetX := cx + i
		targetY := cy

		// clip の範囲内かチェック
		if clip.Contains(targetX, targetY) {
			// キャンバスの物理境界チェック
			if targetY >= 0 && targetY < len(self.model.canvas) &&
				targetX >= 0 && targetX < len(self.model.canvas[targetY]) {
				self.model.canvas[targetY][targetX] = Cell{Rune: r, Color: "white"}
			}
		}
	}
}
