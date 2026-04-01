package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTSpeed struct {
	Base  *yamlui.UISelect
	model *model
}

func NewBTSpeed(m *model, parent *yamlui.UIBase, data map[string]script.Value) *BTSpeed {
	selectUI := &BTSpeed{
		Base:  yamlui.NewUISelect(3),
		model: m,
	}
	selectUI.Base.Base.SetDrawIF(selectUI)
	return selectUI
}

func (self *BTSpeed) Draw(x, y float64, ctx yamlui.DrawContext) {
	x = ctx.ParentClip.AlignCenterX(ctx.Base.W) // 1行あたり16文字分の幅を中央寄せ
	clipX, clipY, clipW, clipH := ctx.Clip.IRect()

	for i, item := range self.Base.Items {
		// 表示領域の高さ(clip.H)を超えたら描画しない
		if i/self.Base.Rows >= clipH {
			break
		}

		// 描画する Y 座標（1行ずつズラしていく）
		y := clipY + i/self.Base.Rows
		x := clipX + (i%self.Base.Rows)*15 // 1行あたり16文字分の幅を確保

		// 選択中の行（SelectNo）なら、カーソルを表示
		prefix := "   "
		if i == int(ctx.Base.SelectNo) {
			prefix = "▶"
			// TODO: ここで Canvas 側に反転色や色の指定を渡せるとリッチになります
		}

		// Canvas への書き込み（1文字ずつ。Label のロジックを再利用）
		line := prefix + item.Base.Text
		for j, char := range line {
			// 横幅(clip.W)を超えないようにガード
			if j >= clipW {
				break
			}
			// Canvas の (x + j, y) に char を書き込む
			// m.Canvas.Set(x + j, y, char)
			self.model.canvas[y][x+j] = Cell{Rune: char, Color: "white"}
		}
	}
}
