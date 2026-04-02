package main

import (
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTTitle struct {
	UIBase *yamlui.UIBase
	model  *model
}

func NewBTTitle(model *model) *BTTitle {
	return &BTTitle{
		UIBase: yamlui.NewUIBase(),
		model:  model,
	}
}

func (self *BTTitle) GetUIBase() *yamlui.UIBase {
	return self.UIBase
}

func (self *BTTitle) Clone() yamlui.UIComponent[*yamlui.UIBase] {
	return &BTTitle{
		UIBase: self.UIBase.Clone(),
		model:  self.model,
	}
}

func (self *BTTitle) Setup(lib *yamlui.YAMLUI, type_ string, parent *yamlui.UIBase, data map[string]script.Value) error {
	if err := self.UIBase.Setup(lib, type_, parent, data); err != nil { // super call
		return err
	}

	self.UIBase.SetDrawIF(self)
	return nil
}

// func stringWidth(s string) int {
// 	width := 0
// 	for _, r := range s {
// 		if r <= 127 {
// 			width += 1
// 		} else {
// 			width += 1
// 		}
// 	}
// 	return width
// }

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
