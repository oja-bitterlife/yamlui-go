package main

import "github.com/oja-bitterlife/yamlui-go/yamlui"

type BTStart struct {
	Base  *yamlui.UISelect
	model *model
}

func NewBTStart(m *model, rows int) *BTSpeed {
	selectUI := &BTSpeed{
		Base:  yamlui.NewUISelect(rows),
		model: m,
	}
	selectUI.Base.Base.SetDrawIF(selectUI)
	return selectUI
}

func (self *BTStart) Draw(ui *yamlui.UIBase, clip yamlui.Area, ctx yamlui.DrawContext) {
	clip.X = ctx.ParentArea.AlignCenterX(clip.W) // 1行あたり16文字分の幅を中央寄せ

	for i, item := range self.Base.Items {
		// 表示領域の高さ(clip.H)を超えたら描画しない
		if i >= clip.H {
			break
		}

		// 描画する Y 座標（1行ずつズラしていく）
		y := clip.Y + i/self.Base.Rows
		x := clip.X + (i % self.Base.Rows)

		// 選択中の行（SelectNo）なら、カーソルを表示
		prefix := "   "
		if i == ui.SelectNo {
			prefix = "▶"
			// TODO: ここで Canvas 側に反転色や色の指定を渡せるとリッチになります
		}

		// Canvas への書き込み（1文字ずつ。Label のロジックを再利用）
		line := prefix + item.Base.Text
		for j, char := range line {
			// 横幅(clip.W)を超えないようにガード
			if j >= clip.W {
				break
			}
			// Canvas の (x + j, y) に char を書き込む
			// m.Canvas.Set(x + j, y, char)
			self.model.canvas[y][x+j] = Cell{Rune: char, Color: "white"}
		}
	}
}
