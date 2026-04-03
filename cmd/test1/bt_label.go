package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTLabel struct {
	UIBase *yamlui.UIBase
	model  *model
}

func NewBTLabel(m *model) *BTLabel {
	return &BTLabel{
		UIBase: yamlui.NewUIBase(),
		model:  m,
	}
}

func (self *BTLabel) GetUIBase() *yamlui.UIBase {
	return self.UIBase
}

func (self *BTLabel) Clone() yamlui.UIComponent[*yamlui.UIBase] {
	return &BTLabel{
		UIBase: self.UIBase.Clone(),
		model:  self.model,
	}
}

func (self *BTLabel) Setup(type_ string, data script.ValueMap) error {
	if err := self.UIBase.Setup(type_, data); err != nil { // super call
		return err
	}

	self.UIBase.SetDrawIF(self)
	return nil
}

func (self *BTLabel) Draw(x, y float64, ctx yamlui.DrawContext) {
	if x < 0 || y < 0 {
		return // 負の座標は描画しない
	}
	if x >= ctx.Clip.Right() || y >= ctx.Clip.Bottom() {
		return // クリップ範囲外は描画しない
	}

	line := self.UIBase.Text
	for j, char := range line {
		// 横幅(clip.W)を超えないようにガード
		if j >= int(ctx.Clip.W) {
			break
		}
		// Canvas の (x + j, y) に char を書き込む
		// m.Canvas.Set(x + j, y, char)
		self.model.canvas[int(y)][int(x)+j] = Cell{Rune: char, Color: "white"}
	}
}
