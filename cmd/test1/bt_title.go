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

func (self *BTTitle) Draw(ui *yamlui.UIBase, clip yamlui.Area, ctx yamlui.DrawContext) {
	// 64x24のエリア内での描画
	logo := []string{
		` __     __      __  __  _      _    _  _____  _ `,
		` \ \   / /     |  \/  || |    | |  | ||_   _|| |`,
		`  \ \_/ /__ _  | \  / || |    | |  | |  | |  | |`,
		`   \   // _` + "`" + ` | | |\/| || |    | |  | |  | |  | |`, // バッククォートのエスケープに注意
		`    | || (_| | | |  | || |____| |__| | _| |_ |_|`,
		`    |_| \__,_| |_|  |_||______|\____/ |_____|(_)`,
	}

	// Y座標：上から 1/4 くらいの場所に置くと DQ っぽいです
	startY := int(clip.Y + 2)

	for i, line := range logo {
		if i >= int(clip.H) {
			break
		}

		// X座標：中央寄せ
		// startX := int(clip.AlignCenterIX(stringWidth(line)) - 1)
		startX := int(clip.X)

		for j, char := range line {
			if startX+j < int(clip.X+clip.W) {
				// Canvasにセット
				self.model.canvas[startY+i][startX+j] = Cell{Rune: char, Color: "white"}
			}
		}
	}
}
