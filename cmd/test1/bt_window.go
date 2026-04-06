package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTWindow struct {
	UIBase *yamlui.UIBase
	model  *model
	orgH   int
}

func NewBTWindow(m *model) *BTWindow {
	return &BTWindow{
		UIBase: yamlui.NewUIBase(),
		model:  m,
	}
}

// **********************************************************************
// UIComponentIFの実装
func (self *BTWindow) GetUIBase() *yamlui.UIBase {
	return self.UIBase
}

func (self *BTWindow) Clone() yamlui.UICloned {
	// modelは共有でOK
	return &BTWindow{
		UIBase: self.GetUIBase().Clone(),
		model:  self.model,
	}
}

func (self *BTWindow) Setup(type_ string, data script.ValueMap) error {
	if err := self.GetUIBase().Setup(type_, data); err != nil { // super call
		return err
	}

	// クローズのアニメーションのために元の高さを保存
	self.orgH = self.GetUIBase().H

	self.GetUIBase().SetUpdateIF(self)
	self.GetUIBase().SetDrawIF(self)
	return nil
}

func (self *BTWindow) Update(lib *yamlui.YAMLUI, events []string) (string, error) {
	// closing中
	if self.GetUIBase().HasProp("close_ratio") {
		// close_ratioに応じて上から閉じていくようにY座標を計算
		self.GetUIBase().H = self.orgH * (100 - self.GetUIBase().PropNum("close_ratio")) / 100
	}

	return "", nil
}

func (self *BTWindow) Draw(x, y int, ctx yamlui.DrawContext) {
	drawArea := ctx.Clip

	// 3. 描画ループ
	for dy := int(drawArea.Top()); dy < int(drawArea.Bottom()); dy++ {
		for dx := int(drawArea.Left()); dx < int(drawArea.Right()); dx++ {

			// dx, dy が本来の自分の領域のどこに当たるかで文字を決める
			isLeft := (dx == x)
			isRight := (dx == x+self.GetUIBase().W-1)
			isTop := (dy == y)
			isBottom := (dy == y+self.orgH-1)

			var r rune
			if isLeft && isTop {
				r = '╔' // 左上
			}
			if isRight && isTop {
				r = '╗' // 右上
			}
			if isLeft && isBottom {
				r = '╚' // 左下
			}
			if isRight && isBottom {
				r = '╝' // 右下
			}
			if r == 0 {
				if isTop || isBottom {
					r = '═' // 上下の線
				} else if isLeft || isRight {
					r = '║' // 左右の線
				} else {
					r = ' ' // 内部はスペース
				}
			}

			// キャンバスの物理境界チェックをして書き込み
			if dy >= 0 && dy < len(self.model.canvas) && dx >= 0 && dx < len(self.model.canvas[dy]) {
				if drawArea.Contains(dx, dy) {
					self.model.canvas[dy][dx] = Cell{Rune: r, Color: "white"}
				}
			}
		}
	}
}
