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
		"ID":     script.NewString(self.ID),
		"Events": script.NewLitList(events),
		"script": self.script.AST,
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
func (self *UIBase) LoadFromValue(value script.Value) {
	for k, v := range value.Map {
		switch k {
		case "ID":
			self.ID = v.Str // Varsに入れず、不変に

		case "script":
			self.script.AST = v

		case "Events":
			if v.Type != script.TypeLitList {
				script.LogErr("Expected Events to be ListType: " + v.Type.String())
				continue
			}

			events := make([]string, len(v.List))
			for i, eventVal := range v.List {
				if eventVal.Type != script.TypeString {
					script.LogErr("Expected Events to be List of String: " + eventVal.Type.String())
					continue
				}
				events[i] = eventVal.Str
			}
			self.Events = events

		// 処理しないもの
		case "children":

		// その他のPropertyは@をつけてVMのVarsに流し込む
		default:
			self.script.Vars["@"+k] = v
		}
	}
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

func (self *UIBase) PropStrDef(key string, def string) string {
	return self.script.Vars.GetStr(key, def)
}

func (self *UIBase) PropNumDef(key string, def int) int {
	return self.script.Vars.GetNum(key, def)
}

func (self *UIBase) PropBoolDef(key string, def bool) bool {
	return self.script.Vars.GetBool(key, def)
}

// ==================================================
// 削除
func (self *UIBase) DeleteProp(key string) {
	delete(self.script.Vars, key)
}

// ==================================================
// Propの存在チェック
func (self *UIBase) HasProp(key string) bool {
	_, ok := self.script.Vars[key]
	return ok
}
