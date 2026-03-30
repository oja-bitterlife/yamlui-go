package yamlui

import (
	"errors"
	"strconv"
)

// **********************************************************************
// 選択アイテム
type UISelectItem struct {
	Base *UIBase
}

func NewUISelectItem(labal string) *UISelectItem {
	item := &UISelectItem{
		Base: NewUIBase("SelectItem"),
	}
	item.Base.Text = labal
	return item
}

func (self *UISelectItem) GetAction() string {
	return self.Base.Action
}

// **********************************************************************
// 選択UI
type UISelect struct {
	Base  *UIBase
	Items []*UISelectItem
	Rows  int
}

func NewUISelect(rows int) *UISelect {
	selectUI := &UISelect{
		Base: NewUIBase("Select"),
		Rows: rows,
	}
	return selectUI
}

func (self *UISelect) GetSelectItem() (*UISelectItem, error) {
	if self.Base.SelectNo < 0 || self.Base.SelectNo >= len(self.Items) {
		return nil, errors.New("select no is out of range: " + strconv.Itoa(self.Base.SelectNo))
	}
	return self.Items[self.Base.SelectNo], nil
}

func (self *UISelect) GetSelectAction() string {
	item, err := self.GetSelectItem()
	if err != nil {
		return ""
	}
	return item.Base.Action
}

func (self *UISelect) GetSelectNo() int {
	return self.Base.SelectNo
}

// ==================================================
// アイテムの追加と削除
func (self *UISelect) AddItem(item *UISelectItem) {
	self.Items = append(self.Items, item)
}

func (self *UISelect) AddItems(items []*UISelectItem) {
	self.Items = append(self.Items, items...)
}

func (self *UISelect) RemoveItem(index int) {
	if index < 0 || index >= len(self.Items) {
		return
	}
	self.Items = append(self.Items[:index], self.Items[index+1:]...)
}

func (self *UISelect) ClearItems() {
	self.Items = nil
}

// ==================================================
// 選択の操作
// step: 移動量（正の数で下/右、負の数で上/左）
// toggle: trueなら端から端へループ、falseなら端で止まる

// 単純な移動（リストで選択肢を配置している場合の移動）
func (self *UISelect) Next(step int, toggle bool) {
	if len(self.Items) == 0 {
		return
	}

	self.Base.SelectNo += step

	if toggle {
		self.Base.SelectNo = (len(self.Items) + self.Base.SelectNo) % len(self.Items)
	} else {
		self.Base.SelectNo = min(max(self.Base.SelectNo, 0), len(self.Items)-1)
	}
}

// グリッドX移動（RowsでItemを折り返している場合の移動）
func (self *UISelect) NextGridX(step int, toggle bool) {
	if len(self.Items) == 0 {
		return
	}

	gridX := self.Base.SelectNo % self.Rows
	gridY := self.Base.SelectNo / self.Rows
	gridX += step

	if toggle {
		gridX = (self.Rows + gridX) % self.Rows
		if gridY*self.Rows+gridX >= len(self.Items) {
			gridX = 0
		}
	} else {
		gridX = min(max(gridX, 0), self.Rows-1)
		if gridY*self.Rows+gridX >= len(self.Items) {
			gridX = (len(self.Items) - 1) % self.Rows
		}
	}

	self.Base.SelectNo = gridY*self.Rows + gridX
}

// グリッドY移動（RowsでItemを折り返している場合の移動）
func (self *UISelect) NextGridY(step int, toggle bool) {
	if len(self.Items) == 0 {
		return
	}

	gridX := self.Base.SelectNo % self.Rows
	gridY := self.Base.SelectNo / self.Rows
	lineNum := (len(self.Items) + self.Rows - 1) / self.Rows
	gridY += step

	if toggle {
		gridY = (lineNum + gridY) % lineNum
	} else {
		gridY = min(max(gridY, 0), lineNum-1)
	}
	if gridY*self.Rows+gridX >= len(self.Items) {
		gridX = (len(self.Items) - 1) % self.Rows
	}

	self.Base.SelectNo = gridY*self.Rows + gridX
}
