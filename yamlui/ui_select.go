package yamlui

type UISelectItem struct {
	Base   *UIBase
	Action string
}

type UISelect struct {
	Base  *UIBase
	Items []*UISelectItem
	GridX int
}

func NewUISelectItem(labal string) *UISelectItem {
	item := &UISelectItem{
		Base: NewUIBase("SelectItem"),
	}
	item.Base.Text = labal
	return item
}

func NewUISelect() *UISelect {
	selectUI := &UISelect{
		Base: NewUIBase("Select"),
	}
	return selectUI
}

func (self *UISelect) Next(step int, toggle bool) {
	if len(self.Items) == 0 {
		return
	}
	self.Base.SelectNo = len(self.Items) + self.Base.SelectNo + step
	if toggle {
		self.Base.SelectNo = self.Base.SelectNo % len(self.Items)
	} else {
		self.Base.SelectNo = min(max(self.Base.SelectNo, 0), len(self.Items)-1)
	}
}

func (self *UISelect) NextGridX(step int, toggle bool) {
	if len(self.Items) == 0 {
		return
	}
	gridX := self.Base.SelectNo % self.GridX
	gridY := self.Base.SelectNo / self.GridX
	gridX = self.GridX + gridX + step
	if toggle {
		gridX = gridX % self.GridX
		if gridY*self.GridX+gridX >= len(self.Items) {
			gridX = 0
		}
	} else {
		gridX = min(gridX, self.GridX-1)
		if gridY*self.GridX+gridX >= len(self.Items) {
			gridX = (len(self.Items) - 1) % self.GridX
		}
	}
	self.Base.SelectNo = gridY*self.GridX + gridX
}

func (self *UISelect) NextGridY(step int, toggle bool) {
	if len(self.Items) == 0 {
		return
	}
	gridX := self.Base.SelectNo % self.GridX
	gridY := self.Base.SelectNo / self.GridX
	lineNum := (len(self.Items) + self.GridX - 1) / self.GridX
	gridY = lineNum + gridY + step
	if toggle {
		gridY = gridY % lineNum
	} else {
		gridY = min(max(gridY, 0), lineNum-1)
	}
	if gridY*self.GridX+gridX >= len(self.Items) {
		gridX = (len(self.Items) - 1) % self.GridX
	}
	self.Base.SelectNo = gridY*self.GridX + gridX
}
