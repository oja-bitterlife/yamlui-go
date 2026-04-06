package yamlui

import (
	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// UIBaseとValueの相互変換
// ==================================================
// UIBaseをValueに変換する関数
func (self *UIBase) ToValue() script.Value {
	// EventsはValueのリストに変換する
	events := make([]script.Value, len(self.Events))
	for i, event := range self.Events {
		events[i] = script.NewString(event)
	}

	return script.NewLitMap(script.ValueMap{
		"ID":          script.NewString(self.ID),
		"UpdateCount": script.NewNumber(self.UpdateCount),
		"Events":      script.NewLitList(events),
		"IsEnable":    script.NewBool(self.IsEnable),
		"Remove":      script.NewBool(self.Remove),
		"X":           script.NewNumber(self.X),
		"Y":           script.NewNumber(self.Y),
		"W":           script.NewNumber(self.W),
		"H":           script.NewNumber(self.H),
		"IsVisible":   script.NewBool(self.IsVisible),
		"Text":        script.NewString(self.Text),
		"SelectNo":    script.NewNumber(self.SelectNo),
		"SelGridX":    script.NewNumber(self.SelGridX),
		"Prop":        script.NewLitMap(self.Prop),
	})
}

// ==================================================
// ValueからUIBaseに流し込む関数
func (self *UIBase) LoadFromValue(value script.Value) error {
	if value.Type != script.TypeLitMap {
		return script.LogErr("Expected Value to be MapType: " + value.Type.String())
	}
	m := value.Map

	// キーがあれば流し込む。なければデフォルト値のまま
	if v, ok := m["ID"]; ok {
		self.ID = v.Str
	}
	if v, ok := m["UpdateCount"]; ok {
		self.UpdateCount = int(v.Num)
	}

	// EventsはStringのリストであることを確認してから流し込む
	if v, ok := m["Events"]; ok {
		if v.Type != script.TypeLitList {
			return script.LogErr("Expected Events to be ListType: " + v.Type.String())
		}

		events := make([]string, len(v.List))
		for i, eventVal := range v.List {
			if eventVal.Type != script.TypeString {
				return script.LogErr("Expected Events to be List of String: " + eventVal.Type.String())
			}
			events[i] = eventVal.Str
		}
		self.Events = events
	}

	if v, ok := m["IsEnable"]; ok {
		self.IsEnable = v.Bool

	}
	if v, ok := m["Remove"]; ok {
		self.Remove = v.Bool
	}
	if v, ok := m["X"]; ok {
		self.X = v.Num
	}
	if v, ok := m["Y"]; ok {
		self.Y = v.Num
	}
	if v, ok := m["W"]; ok {
		self.W = v.Num
	}
	if v, ok := m["H"]; ok {
		self.H = v.Num
	}
	if v, ok := m["IsVisible"]; ok {
		self.IsVisible = v.Bool
	}
	if v, ok := m["Text"]; ok {
		self.Text = v.Str
	}
	if v, ok := m["SelectNo"]; ok {
		self.SelectNo = v.Num
	}
	if v, ok := m["SelGridX"]; ok {
		self.SelGridX = v.Num
	}

	// scriptはVMの作成とcmdの登録を行う
	if v, ok := m["script"]; ok {
		script, err := script.NewVM(v)
		if err != nil {
			return err
		} else {
			self.setUIScriptCmds(script)
			self.script = script
		}
	}

	// PropはMapTypeであればそのまま流し込む。なければ空のまま
	if prop, ok := m["Prop"]; ok {
		if prop.Type != script.TypeLitMap {
			return script.LogErr("Expected Prop to be MapType: " + prop.Type.String())
		}
	}

	return nil
}

// **********************************************************************
// Propのgetter/setter
// 設定
func (self *UIBase) SetPropStr(key string, str string) {
	self.Prop[key] = script.NewString(str)
}

func (self *UIBase) AddPropNum(key string, num int) {
	self.Prop[key] = script.NewNumber(self.PropNum(key) + num)
}

func (self *UIBase) SetPropNum(key string, num int) {
	self.Prop[key] = script.NewNumber(num)
}

func (self *UIBase) SetPropBool(key string, b bool) {
	self.Prop[key] = script.NewBool(b)
}

// 取り出すだけ
func (self *UIBase) PropStr(key string) string {
	return self.Prop[key].Str
}

func (self *UIBase) PropNum(key string) int {
	return self.Prop[key].Num
}

func (self *UIBase) PropBool(key string) bool {
	return self.Prop[key].Bool
}

// 取り出してクリアする
func (self *UIBase) TakePropStr(key string) string {
	str := self.Prop[key].Str
	delete(self.Prop, key)
	return str
}

func (self *UIBase) TakePropNum(key string) int {
	num := self.Prop[key].Num
	delete(self.Prop, key)
	return num
}

func (self *UIBase) TakePropBool(key string) bool {
	b := self.Prop[key].Bool
	delete(self.Prop, key)
	return b
}
