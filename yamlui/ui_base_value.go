package yamlui

import (
	"errors"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// UIBaseとValueの相互変換
// ==================================================
// UIBaseをValueに変換する関数
func (self *UIBase) ToValue() script.Value {
	return script.NewLitMap(map[string]script.Value{
		"Type":      script.NewString(self.Type),
		"ID":        script.NewString(self.ID),
		"IsEnable":  script.NewBool(self.IsEnable),
		"X":         script.NewNumber(float64(self.X)),
		"Y":         script.NewNumber(float64(self.Y)),
		"W":         script.NewNumber(float64(self.W)),
		"H":         script.NewNumber(float64(self.H)),
		"IsVisible": script.NewBool(self.IsVisible),
		"Text":      script.NewString(self.Text),
		"Color":     script.NewString(self.Color),
		"SelectNo":  script.NewNumber(float64(self.SelectNo)),
		"SelGridX":  script.NewNumber(float64(self.SelGridX)),
		"Return":    script.NewString(self.ScriptReturn),
		"Prop":      script.NewLitMap(self.Prop),
	})
}

// ==================================================
// ValueからUIBaseに流し込む関数
func (self *UIBase) LoadFromValue(value script.Value) error {
	if value.Type != script.TypeLitMap {
		return errors.New("Expected Value to be MapType: " + value.Type.String())
	}
	m := value.Map

	// キーがあれば流し込む。なければデフォルト値のまま
	if v, ok := m["Type"]; ok {
		self.Type = v.Str
	}
	if v, ok := m["ID"]; ok {
		self.ID = v.Str
	}
	if v, ok := m["IsEnable"]; ok {
		self.IsEnable = v.Bool
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
	if v, ok := m["Color"]; ok {
		self.Color = v.Str
	}
	if v, ok := m["SelectNo"]; ok {
		self.SelectNo = v.Num
	}
	if v, ok := m["SelGridX"]; ok {
		self.SelGridX = v.Num
	}
	if v, ok := m["Return"]; ok {
		self.ScriptReturn = v.Str
	}

	if prop, ok := m["Prop"]; ok {
		if prop.Type != script.TypeLitMap {
			return errors.New("Expected Prop to be MapType: " + prop.Type.String())
		}
	}
	return nil
}
