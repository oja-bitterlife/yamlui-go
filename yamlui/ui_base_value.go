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

	// UIBaseの基本的なプロパティをValueMapに入れる
	data := script.ValueMap{
		"ID":       script.NewString(self.ID),
		"Events":   script.NewLitList(events),
		"IsEnable": script.NewBool(self.IsEnable),
		"Remove":   script.NewBool(self.Remove),
		"X":        script.NewNumber(self.X),
		"Y":        script.NewNumber(self.Y),
		"W":        script.NewNumber(self.W),
		"H":        script.NewNumber(self.H),
		"script":   self.script.AST,
	}

	// Propertyは@付きの名前でVarsに入っているものをValueMapに入れる
	for k, v := range self.script.Vars {
		// @付きだけ回収する.@がないものはlocal変数
		if len(k) > 0 && k[0] == '@' {
			// @を外して入れる
			data[k[1:]] = v
		}
	}

	return script.NewLitMap(data)
}

// ==================================================
// ValueからUIBaseに流し込む関数
func (self *UIBase) LoadFromValue(value script.Value) error {
	if value.Type != script.TypeLitMap {
		return script.LogErr("Expected Value to be MapType: " + value.Type.String())
	}
	m := value.Map

	for k, v := range m {
		switch k {
		case "ID":
			self.ID = v.Str
		case "Events":
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
		case "IsEnable":
			self.IsEnable = v.Bool()
		case "Remove":
			self.Remove = v.Bool()
		case "X":
			self.X = v.Num
		case "Y":
			self.Y = v.Num
		case "W":
			self.W = v.Num
		case "H":
			self.H = v.Num
		case "script":
			self.script.AST = v

		case "Type", "children":
			// 処理しないもの
		default:
			// その他のPropertyは@をつけてVMのVarsに流し込む
			self.script.Vars["@"+k] = v
		}
	}

	return nil
}

// **********************************************************************
// Propのgetter/setter
// ==================================================
// 設定
func (self *UIBase) SetPropStr(key string, str string) {
	self.script.Vars[key] = script.NewString(str)
}

func (self *UIBase) AddPropNum(key string, num int) {
	self.script.Vars[key] = script.NewNumber(self.PropNum(key) + num)
}

func (self *UIBase) SetPropNum(key string, num int) {
	self.script.Vars[key] = script.NewNumber(num)
}

func (self *UIBase) SetPropBool(key string, b bool) {
	self.script.Vars[key] = script.NewBool(b)
}

// ==================================================
// 取り出すだけ
func (self *UIBase) PropStr(key string) string {
	return self.script.Vars.GetStr(key, "")
}

func (self *UIBase) PropNum(key string) int {
	return self.script.Vars.GetNum(key, 0)
}

func (self *UIBase) PropBool(key string) bool {
	return self.script.Vars.GetBool(key, false)
}

// ==================================================
// 取り出してクリアする
func (self *UIBase) TakePropStr(key string) string {
	str := self.PropStr(key)
	delete(self.script.Vars, key)
	return str
}

func (self *UIBase) TakePropNum(key string) int {
	num := self.PropNum(key)
	delete(self.script.Vars, key)
	return num
}

func (self *UIBase) TakePropBool(key string) bool {
	b := self.PropBool(key)
	delete(self.script.Vars, key)
	return b
}

// ==================================================
// Propの存在チェック
func (self *UIBase) HasProp(key string) bool {
	_, ok := self.script.Vars[key]
	return ok
}
