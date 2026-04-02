package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTWindow struct {
	WinBase *yamlui.UIWindow
	model   *model
}

func NewBTWindow(m *model) *BTWindow {
	return &BTWindow{
		WinBase: yamlui.NewUIWindow(),
		model:   m,
	}
}

// **********************************************************************
// UIComponentIFの実装
func (self *BTWindow) GetUIBase() *yamlui.UIBase {
	return self.WinBase.UIBase
}

func (self *BTWindow) Clone() yamlui.UIComponent[*yamlui.UIBase] {
	// modelは共有でOK
	return &BTWindow{
		WinBase: self.WinBase.Clone().(*yamlui.UIWindow),
		model:   self.model,
	}
}

func (self *BTWindow) Setup(lib *yamlui.YAMLUI, type_ string, parent *yamlui.UIBase, data map[string]script.Value) error {
	self.WinBase.UIBase.SetDrawIF(self)
	return nil
}

func (self *BTWindow) Draw(x, y float64, ctx yamlui.DrawContext) {
	drawArea := ctx.Clip
	myArea := yamlui.Area{
		X: x,
		Y: y,
		W: ctx.Clip.W,
		H: ctx.Clip.H,
	}

	// 3. 描画ループ
	for dy := int(drawArea.Top()); dy < int(drawArea.Bottom()); dy++ {
		for dx := int(drawArea.Left()); dx < int(drawArea.Right()); dx++ {

			// 現在の (x, y) が「本来の自分の領域(myArea)」のどこに当たるかで文字を決める
			isLeft := (dx == myArea.ILeft())
			isRight := (dx == myArea.IRight()-1)
			isTop := (dy == myArea.ITop())
			isBottom := (dy == myArea.IBottom()-1)

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
				self.model.canvas[dy][dx] = Cell{Rune: r, Color: "white"}
			}
		}
	}
}
