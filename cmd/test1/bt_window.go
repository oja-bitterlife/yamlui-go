package main

import "github.com/oja-bitterlife/yamlui-go/yamlui"

type BTWindow struct {
	Base  *yamlui.UIWindow
	model *model
}

func NewBTWindow(m *model) *BTWindow {
	window := &BTWindow{
		Base:  yamlui.NewUIWindow(),
		model: m,
	}
	window.Base.Base.SetDrawIF(window)
	return window
}
func NewBTWindowRect(m *model, x, y, w, h float64) *BTWindow {
	window := &BTWindow{
		Base:  yamlui.NewUIWindow(),
		model: m,
	}
	window.Base.Base.SetRect(x, y, w, h)
	window.Base.Base.SetDrawIF(window)
	return window
}

func (w *BTWindow) Remap(ui *yamlui.UIBase, data map[string]any) (*yamlui.UIBase, error) {
	w.Base.Base.CopyProp(ui)
	return w.Base.Base, nil
}

func (w *BTWindow) Draw(ui *yamlui.UIBase, clip yamlui.Area, ctx yamlui.DrawContext) {
	// 1. 自分の本来の絶対座標領域を取得
	myArea := ui.Area() // DrawTree で ui.X, ui.Y が絶対座標に更新されている前提

	// 2. 親の制限(clip)と自分の領域(myArea)の重なり＝「実際に描画して良い範囲」
	drawArea := myArea.Clip(clip)

	// 3. 描画ループ
	for y := int(drawArea.Top()); y < int(drawArea.Bottom()); y++ {
		for x := int(drawArea.Left()); x < int(drawArea.Right()); x++ {

			// 現在の (x, y) が「本来の自分の領域(myArea)」のどこに当たるかで文字を決める
			isLeft := (x == myArea.ILeft())
			isRight := (x == myArea.IRight()-1)
			isTop := (y == myArea.ITop())
			isBottom := (y == myArea.IBottom()-1)

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
			if y >= 0 && y < len(w.model.canvas) && x >= 0 && x < len(w.model.canvas[y]) {
				w.model.canvas[y][x] = Cell{Rune: r, Color: "white"}
			}
		}
	}
}
