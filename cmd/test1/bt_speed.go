package main

import (
	"strings"

	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTSpeed struct {
	Base  *yamlui.UISelect
	model *model
	texts []string
}

func NewBTSpeed(m *model, parent *yamlui.UIBase, data map[string]script.Value) *BTSpeed {
	// ,区切りのTextを分解してTrimしてtextsに格納
	texts := strings.Split(data["Text"].Str, ",")
	for i := range texts {
		texts[i] = strings.TrimSpace(texts[i])
	}

	selectUI := &BTSpeed{
		Base:  yamlui.NewUISelect(len(texts), len(texts)),
		model: m,
		texts: texts,
	}
	selectUI.Base.Base.SetDrawIF(selectUI)
	return selectUI
}

func (self *BTSpeed) Draw(x, y float64, ctx yamlui.DrawContext) {
	clipW, clipH := ctx.Clip.IWH()

	for i := range self.Base.ItemNum {
		// 表示領域の高さ(clip.H)を超えたら描画しない
		if i/self.Base.RowNum >= clipH {
			break
		}

		// 描画する Y 座標（1行ずつズラしていく）
		drawX := x + float64(i%self.Base.RowNum)*15 // 1行あたり16文字分の幅を確保
		drawY := y + float64(i/self.Base.RowNum)

		// 選択中の行（SelectNo）なら、カーソルを表示
		prefix := "   "
		if i == int(ctx.Base.SelectNo) {
			prefix = "▶"
			// TODO: ここで Canvas 側に反転色や色の指定を渡せるとリッチになります
		}

		// Canvas への書き込み（1文字ずつ。Label のロジックを再利用）
		line := prefix + self.texts[i]
		for j, char := range line {
			// 横幅(clip.W)を超えないようにガード
			if j >= clipW {
				break
			}
			// Canvas の (x + j, y) に char を書き込む
			// m.Canvas.Set(x + j, y, char)
			self.model.canvas[int(drawY)][int(drawX)+j] = Cell{Rune: char, Color: "white"}
		}
	}
}
