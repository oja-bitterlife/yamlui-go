package main

import (
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTWindow struct {
	UIBase *yamlui.UIBase
	model  *model
	orgH   int
}

func NewBTWindow(lib *yamlui.YAMLUI, m *model) *BTWindow {
	return &BTWindow{
		UIBase: yamlui.NewUIBase(lib),
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

func (self *BTWindow) Setup(type_ string) {
	self.GetUIBase().Setup(type_)

	// クローズのアニメーションのために元の高さを保存
	uiArea := self.GetUIBase().Area()
	self.orgH = uiArea.H

	self.GetUIBase().SetDrawIF(self)
}

func (self *BTWindow) Draw(lib *yamlui.YAMLUI, x, y int, clip yamlui.Area) {
	drawArea := clip
	uiArea := self.GetUIBase().Area()

	// closing中
	if self.GetUIBase().HasProp("close_ratio") {
		// close_ratioに応じて上から閉じていくようにY座標を計算
		drawArea.H = self.orgH * (100 - self.GetUIBase().PropNum("close_ratio")) / 100
	}

	// 3. 描画ループ
	for dy := int(drawArea.Top()); dy < int(drawArea.Bottom()); dy++ {
		for dx := int(drawArea.Left()); dx < int(drawArea.Right()); dx++ {

			// dx, dy が本来の自分の領域のどこに当たるかで文字を決める
			isLeft := (dx == x)
			isRight := (dx == x+uiArea.W-1)
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
					if r != ' ' {
						self.model.canvas[dy][dx] = Cell{Rune: r, Color: "white"}
					}
				}
			}
		}
	}
}
