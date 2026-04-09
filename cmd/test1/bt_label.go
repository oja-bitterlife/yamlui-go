package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTLabel struct {
	UIBase *yamlui.UIBase
	model  *model
}

func NewBTLabel(lib *yamlui.YAMLUI, m *model) *BTLabel {
	return &BTLabel{
		UIBase: yamlui.NewUIBase(lib),
		model:  m,
	}
}

func (self *BTLabel) GetUIBase() *yamlui.UIBase {
	return self.UIBase
}

func (self *BTLabel) Clone() yamlui.UICloned {
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

func (self *BTLabel) Draw(lib *yamlui.YAMLUI, x, y int, clip yamlui.Area) {
	if x < 0 || y < 0 {
		return // 負の座標は描画しない
	}
	if x >= clip.Right() || y >= clip.Bottom() {
		return // クリップ範囲外は描画しない
	}

	line := self.GetUIBase().PropStr("@Text")
	for j, char := range line {
		// 横幅(clip.W)を超えないようにガード
		if j >= int(clip.W) {
			break
		}
		// Canvas の (x + j, y) に char を書き込む
		// m.Canvas.Set(x + j, y, char)
		self.model.canvas[int(y)][int(x)+j] = Cell{Rune: char, Color: "white"}
	}
}
