package main

import (
	"strings"

	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTStart struct {
	SelBase *yamlui.UISelect
	model   *model
	texts   []string
}

func NewBTStart(lib *yamlui.YAMLUI, model *model) *BTStart {
	start := &BTStart{
		SelBase: yamlui.NewUISelect(lib, 0, 1),
		model:   model,
	}
	start.GetUIBase().SetDrawIF(start)
	start.GetUIBase().SetDispatchIF(start)
	return start
}

func (self *BTStart) GetUIBase() *yamlui.UIBase {
	return self.SelBase.UIBase
}

func (self *BTStart) Clone() yamlui.UICloned {
	return &BTStart{
		SelBase: self.SelBase.Clone().(*yamlui.UISelect),
		model:   self.model,
		texts:   self.texts,
	}
}

// Dispatch
func (self *BTStart) Dispatch(lib *yamlui.YAMLUI, event string) {
	switch event {
	case yamlui.EVENT_UI_ONCREATE:
		lib.Log(self.GetUIBase().PropStr(yamlui.PROP_TEXT))
		// ,区切りのTextを分解してTrimしてtextsに格納
		texts := strings.Split(self.GetUIBase().PropStr(yamlui.PROP_TEXT), ",")
		for i := range texts {
			texts[i] = strings.TrimSpace(texts[i])
		}
		self.texts = texts

		// ItemNumをtextsの数に設定
		self.SelBase.ItemNum = len(texts)
	case "key:up":
		self.SelBase.NextGridY(-1, true)

	case "key:down":
		self.SelBase.NextGridY(1, true)
	}
}

// Draw
func (self *BTStart) Draw(lib *yamlui.YAMLUI, x, y int, clip yamlui.Area) {
	for i := range self.SelBase.ItemNum {
		// 表示領域の高さ(clip.H)を超えたら描画しない
		if i >= int(clip.H) {
			break
		}

		// 描画する Y 座標（1行ずつズラしていく）
		y := clip.Y + i/self.SelBase.RowsNum
		x := clip.X + (i % self.SelBase.RowsNum)

		// 選択中の行（SelectNo）なら、カーソルを表示
		line := "   "
		if i == int(self.SelBase.SelectNo()) {
			line = "▶"
			// TODO: ここで Canvas 側に反転色や色の指定を渡せるとリッチになります
		}

		// Canvas への書き込み（1文字ずつ。Label のロジックを再利用）
		line += self.texts[i]
		for j, char := range line {
			// 横幅(clip.W)を超えないようにガード
			if j >= int(clip.W) {
				break
			}
			// Canvas の (x + j, y) に char を書き込む
			// m.Canvas.Set(x + j, y, char)
			self.model.canvas[y][x+j] = Cell{Rune: char, Color: "white"}
		}
	}
}
