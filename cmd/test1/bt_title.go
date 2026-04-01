package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTTitle struct {
	Base  *yamlui.UILabel
	model *model
}

func NewBTTitle(m *model, ui *yamlui.UIBase, data map[string]script.Value) *BTTitle {
	title := &BTTitle{
		Base:  yamlui.NewUILabel("YAMLUI"),
		model: m,
	}
	title.Base.Base.SetDrawIF(title)
	return title
}

func stringWidth(s string) int {
	width := 0
	for _, r := range s {
		if r <= 127 {
			width += 1
		} else {
			width += 1
		}
	}
	return width
}

func (self *BTTitle) Draw(x, y float64, ctx yamlui.DrawContext) {
	// 64x24のエリア内での描画
	logo := []string{
		` __     __      __  __  _      _    _  _____  _ `,
		` \ \   / /     |  \/  || |    | |  | ||_   _|| |`,
		`  \ \_/ /__ _  | \  / || |    | |  | |  | |  | |`,
		`   \   // _` + "`" + ` | | |\/| || |    | |  | |  | |  | |`, // バッククォートのエスケープに注意
		`    | || (_| | | |  | || |____| |__| | _| |_ |_|`,
		`    |_| \__,_| |_|  |_||______|\____/ |_____|(_)`,
	}

	for i, line := range logo {
		drawY := int(y) + i
		if ctx.Clip.IContainsY(drawY) == false {
			continue
		}

		for j, char := range line {
			drawX := int(x) + j
			if ctx.Clip.IContainsX(drawX) == false {
				continue
			}
			self.model.canvas[drawY][drawX] = Cell{Rune: char, Color: "white"}
		}
	}
}
