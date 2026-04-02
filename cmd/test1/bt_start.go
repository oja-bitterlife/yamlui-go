package main

import (
	"strings"

	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTStart struct {
	SelBase *yamlui.UISelect
	model   *model
	texts   []string
}

func NewBTStart(model *model) *BTStart {
	return &BTStart{
		SelBase: yamlui.NewUISelect(0, 1), // ItemNumは後でSetupで設定する
		model:   model,
	}
}

func (self *BTStart) GetUIBase() *yamlui.UIBase {
	return self.SelBase.UIBase
}

func (self *BTStart) Clone() yamlui.UIComponent[*yamlui.UIBase] {
	return &BTStart{
		SelBase: self.SelBase.Clone().(*yamlui.UISelect),
		model:   self.model,
		texts:   self.texts,
	}
}

func (self *BTStart) Setup(lib *yamlui.YAMLUI, type_ string, parent *yamlui.UIBase, data map[string]script.Value) error {
	if err := self.SelBase.Setup(lib, type_, parent, data); err != nil { // super call
		return err
	}

	// ,区切りのTextを分解してTrimしてtextsに格納
	texts := strings.Split(data["Text"].Str, ",")
	for i := range texts {
		texts[i] = strings.TrimSpace(texts[i])
	}
	self.texts = texts

	// ItemNumをtextsの数に設定
	self.SelBase.ItemNum = len(texts)
	self.SelBase.RowsNum = max(int(data["RowsNum"].Num), 1) // RowsNumは1以上にする

	self.SelBase.UIBase.SetDrawIF(self)
	self.SelBase.UIBase.SetUpdateIF(self)

	return nil
}

// Update
func (self *BTStart) Update(ctx yamlui.UpdateContext) error {
	for _, event := range ctx.Events {
		if event == "key:up" {
			self.SelBase.NextGridY(-1, true)
		}
		if event == "key:down" {
			self.SelBase.NextGridY(1, true)
		}
		if event == "key:enter" {
			self.GetUIBase().Action = self.texts[int(self.GetUIBase().SelectNo)]
		}
	}
	return nil
}

// Draw
func (self *BTStart) Draw(x, y float64, ctx yamlui.DrawContext) {

	for i := range self.SelBase.ItemNum {
		// 表示領域の高さ(clip.H)を超えたら描画しない
		if i >= int(ctx.Clip.H) {
			break
		}

		// 描画する Y 座標（1行ずつズラしていく）
		y := ctx.Clip.IY() + i/self.SelBase.RowsNum
		x := ctx.Clip.IX() + (i % self.SelBase.RowsNum)

		// 選択中の行（SelectNo）なら、カーソルを表示
		line := "   "
		if i == int(self.GetUIBase().SelectNo) {
			line = "▶"
			// TODO: ここで Canvas 側に反転色や色の指定を渡せるとリッチになります
		}

		// Canvas への書き込み（1文字ずつ。Label のロジックを再利用）
		line += self.texts[i]
		for j, char := range line {
			// 横幅(clip.W)を超えないようにガード
			if j >= int(ctx.Clip.W) {
				break
			}
			// Canvas の (x + j, y) に char を書き込む
			// m.Canvas.Set(x + j, y, char)
			self.model.canvas[y][x+j] = Cell{Rune: char, Color: "white"}
		}
	}
}
