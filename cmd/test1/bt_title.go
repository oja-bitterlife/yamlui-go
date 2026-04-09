package main

import (
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type BTTitle struct {
	UIBase *yamlui.UIBase
	model  *model
}

func NewBTTitle(lib *yamlui.YAMLUI, model *model) *BTTitle {
	return &BTTitle{
		UIBase: yamlui.NewUIBase(lib),
		model:  model,
	}
}

func (self *BTTitle) GetUIBase() *yamlui.UIBase {
	return self.UIBase
}

func (self *BTTitle) Clone() yamlui.UICloned {
	return &BTTitle{
		UIBase: self.UIBase.Clone(),
		model:  self.model,
	}
}

func (self *BTTitle) Setup(type_ string) {
	self.UIBase.Setup(type_)
	self.UIBase.SetDrawIF(self)
}

func (self *BTTitle) Draw(lib *yamlui.YAMLUI, x, y int, clip yamlui.Area) {
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
		if clip.ContainsY(drawY) == false {
			continue
		}

		for j, char := range line {
			drawX := int(x) + j
			if clip.ContainsX(drawX) == false {
				continue
			}
			self.model.canvas[drawY][drawX] = Cell{Rune: char, Color: "white"}
		}
	}
}
