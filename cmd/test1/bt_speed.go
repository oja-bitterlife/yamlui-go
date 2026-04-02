package main

import (
	"strings"

	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTSpeed struct {
	SelBase *yamlui.UISelect
	model   *model
	texts   []string
}

func NewBTSpeed(model *model) *BTSpeed {
	return &BTSpeed{
		SelBase: yamlui.NewUISelect(0, 1), // ItemNumは後でSetupで設定する
		model:   model,
	}
}

func (self *BTSpeed) GetUIBase() *yamlui.UIBase {
	return self.SelBase.UIBase
}

func (self *BTSpeed) Clone() yamlui.UIComponent[*yamlui.UIBase] {
	return &BTSpeed{
		SelBase: self.SelBase.Clone().(*yamlui.UISelect),
		model:   self.model,
		texts:   self.texts,
	}
}

func (self *BTSpeed) Setup(lib *yamlui.YAMLUI, type_ string, parent *yamlui.UIBase, data map[string]script.Value) error {
	// ,区切りのTextを分解してTrimしてtextsに格納
	texts := strings.Split(data["Text"].Str, ",")
	for i := range texts {
		texts[i] = strings.TrimSpace(texts[i])
	}
	self.texts = texts

	// ItemNumをtextsの数に設定
	self.SelBase.ItemNum = len(texts)
	self.SelBase.RowNum = max(int(data["RowNum"].Num), 1) // RowNumは1以上にする

	self.SelBase.UIBase.SetUpdateIF(self)
	self.SelBase.UIBase.SetDrawIF(self)

	return nil
}

// Update
func (self *BTSpeed) Update(ctx yamlui.UpdateContext) error {
	for _, event := range ctx.Events {
		if event == "key:left" {
			self.SelBase.NextGridX(-1, true)
		}
		if event == "key:right" {
			self.SelBase.NextGridX(1, true)
		}
	}
	return nil
}

func (self *BTSpeed) Draw(x, y float64, ctx yamlui.DrawContext) {
	clipW, clipH := ctx.Clip.IWH()

	for i := range self.SelBase.ItemNum {
		// 表示領域の高さ(clip.H)を超えたら描画しない
		if i/self.SelBase.RowNum >= clipH {
			break
		}

		// 描画する Y 座標（1行ずつズラしていく）
		drawX := x + float64(i%self.SelBase.RowNum)*15 // 1行あたり16文字分の幅を確保
		drawY := y + float64(i/self.SelBase.RowNum)

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
