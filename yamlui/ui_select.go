package yamlui

import "github.com/oja-bitterlife/yamlui-go/script"

// **********************************************************************
// 選択UI
type UISelect struct {
	UIBase  *UIBase
	ItemNum int
	RowNum  int
}

func NewUISelect(itemNum, rowNum int) *UISelect {
	return &UISelect{
		UIBase:  NewUIBase(),
		ItemNum: itemNum,
		RowNum:  rowNum,
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

func (self *UISelect) Clone() *UISelect {
	return &UISelect{
		UIBase:  self.UIBase.Clone(),
		ItemNum: self.ItemNum,
		RowNum:  self.RowNum,
	}
}

func (self *UISelect) Setup(lib *YAMLUI, type_ string, parent *UIBase, data map[string]script.Value) error {
	return nil
}

// **********************************************************************
// 選択の操作
// step: 移動量（正の数で下/右、負の数で上/左）
// toggle: trueなら端から端へループ、falseなら端で止まる
// 単純な移動（リストで選択肢を配置している場合の移動）
func (self *UISelect) Next(step int, toggle bool) {
	if self.RowNum == 0 { // 0で割るのを防止
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
	if self.RowNum == 0 { // 0で割るのを防止
		return
	}

	gridX := self.GetSelectNo() % self.RowNum
	gridY := self.GetSelectNo() / self.RowNum
	gridX += step

	if toggle {
		gridX = (self.RowNum + gridX) % self.RowNum
		if gridY*self.RowNum+gridX >= self.ItemNum {
			gridX = 0
		}
	} else {
		gridX = min(max(gridX, 0), self.RowNum-1)
		if gridY*self.RowNum+gridX >= self.ItemNum {
			gridX = (self.ItemNum - 1) % self.RowNum
		}
	}

	self.SetSelectNo(gridY*self.RowNum + gridX)
}

// グリッドY移動（RowsでItemを折り返している場合の移動）
func (self *UISelect) NextGridY(step int, toggle bool) {
	if self.RowNum == 0 { // 0で割るのを防止
		return
	}

	gridX := self.GetSelectNo() % self.RowNum
	gridY := self.GetSelectNo() / self.RowNum
	lineNum := (self.ItemNum + self.RowNum - 1) / self.RowNum
	gridY += step

	if toggle {
		gridY = (lineNum + gridY) % lineNum
	} else {
		gridY = min(max(gridY, 0), lineNum-1)
	}
	if gridY*self.RowNum+gridX >= self.ItemNum {
		gridX = (self.ItemNum - 1) % self.RowNum
	}

	self.SetSelectNo(gridY*self.RowNum + gridX)
}
