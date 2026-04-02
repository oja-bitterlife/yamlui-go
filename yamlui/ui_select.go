package yamlui

import "github.com/oja-bitterlife/yamlui-go/script"

// **********************************************************************
// 選択UI
type UISelect struct {
	UIBase  *UIBase
	ItemNum int
	RowsNum int
}

func NewUISelect(itemNum, rowsNum int) *UISelect {
	return &UISelect{
		UIBase:  NewUIBase(),
		ItemNum: itemNum,
		RowsNum: rowsNum,
	}
}

func (self *UISelect) GetSelectNo() int {
	return int(self.UIBase.SelectNo)
}

func (self *UISelect) SetSelectNo(selectNo int) {
	self.UIBase.SelectNo = float64(selectNo)
}

// **********************************************************************
// UIComponentIFの実装
func (self *UISelect) GetUIBase() *UIBase {
	return self.UIBase
}

func (self *UISelect) Clone() UIComponent[*UIBase] {
	return &UISelect{
		UIBase:  self.UIBase.Clone(),
		ItemNum: self.ItemNum,
		RowsNum: self.RowsNum,
	}
}

func (self *UISelect) Setup(lib *YAMLUI, type_ string, parent *UIBase, data map[string]script.Value) error {
	return self.UIBase.Setup(lib, type_, parent, data) // super call
}

// **********************************************************************
// 選択の操作
// step: 移動量（正の数で下/右、負の数で上/左）
// toggle: trueなら端から端へループ、falseなら端で止まる
// 単純な移動（リストで選択肢を配置している場合の移動）
func (self *UISelect) Next(step int, toggle bool) {
	if self.RowsNum == 0 { // 0で割るのを防止
		return
	}

	selectNo := self.GetSelectNo() + step

	if toggle {
		selectNo = (self.ItemNum + selectNo) % self.ItemNum
	} else {
		selectNo = min(max(selectNo, 0), self.ItemNum-1)
	}

	self.SetSelectNo(selectNo)
}

// グリッドX移動（RowsでItemを折り返している場合の移動）
func (self *UISelect) NextGridX(step int, toggle bool) {
	if self.RowsNum == 0 { // 0で割るのを防止
		return
	}

	gridX := self.GetSelectNo() % self.RowsNum
	gridY := self.GetSelectNo() / self.RowsNum
	gridX += step

	if toggle {
		gridX = (self.RowsNum + gridX) % self.RowsNum
		if gridY*self.RowsNum+gridX >= self.ItemNum {
			gridX = 0
		}
	} else {
		gridX = min(max(gridX, 0), self.RowsNum-1)
		if gridY*self.RowsNum+gridX >= self.ItemNum {
			gridX = (self.ItemNum - 1) % self.RowsNum
		}
	}

	self.SetSelectNo(gridY*self.RowsNum + gridX)
}

// グリッドY移動（RowsでItemを折り返している場合の移動）
func (self *UISelect) NextGridY(step int, toggle bool) {
	if self.RowsNum == 0 { // 0で割るのを防止
		return
	}

	gridX := self.GetSelectNo() % self.RowsNum
	gridY := self.GetSelectNo() / self.RowsNum
	lineNum := (self.ItemNum + self.RowsNum - 1) / self.RowsNum
	gridY += step

	if toggle {
		gridY = (lineNum + gridY) % lineNum
	} else {
		gridY = min(max(gridY, 0), lineNum-1)
	}
	if gridY*self.RowsNum+gridX >= self.ItemNum {
		gridX = (self.ItemNum - 1) % self.RowsNum
	}

	self.SetSelectNo(gridY*self.RowsNum + gridX)
}
