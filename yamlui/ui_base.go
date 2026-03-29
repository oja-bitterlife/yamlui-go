package yamlui

import "github.com/oja-bitterlife/yamlui-go/script"

// **********************************************************************
// UIのタイプ
type UIType int

const (
	UITypeNone UIType = iota
	UITypeLabel
	UITypeButton
	UITypeSelectH
	UITypeSelectV
	UITypeSelectGrid
	UITypeSelectItem
	UITypeWindow
	UITypeInput
	UITypeDial
)

func (t UIType) String() string {
	switch t {
	case UITypeNone:
		return "None"
	case UITypeLabel:
		return "Label"
	case UITypeButton:
		return "Button"
	case UITypeSelectH:
		return "SelectH"
	case UITypeSelectV:
		return "SelectV"
	case UITypeSelectGrid:
		return "SelectGrid"
	case UITypeSelectItem:
		return "SelectItem"
	case UITypeWindow:
		return "Window"
	case UITypeInput:
		return "Input"
	case UITypeDial:
		return "Dial"
	default:
		return "Unknown"
	}
}

// **********************************************************************
// UIの基本構造.これを保有して各UI構造体を作る
type UIBase struct {
	// 保存されるプロパティ
	Type      UIType
	Count     int
	IsEnable  bool
	IsVisivle bool
	IsAbs     bool
	X         int
	Y         int
	Width     int
	Height    int
	SelectNo  int
	Action    string
	SelGridX  int // SelectGridで横の折り返し位置

	// 保存しないもの
	script *script.Runtime
}

// ==================================================
// VMとのやりとり
func (ui *UIBase) StoreToVM(vm *script.VM) {
	vm.SetVar("@Count", script.NewNumber(float64(ui.Count)))
	vm.SetVar("@IsEnable", script.NewBool(ui.IsEnable))
	vm.SetVar("@IsVisivle", script.NewBool(ui.IsVisivle))
	vm.SetVar("@IsAbs", script.NewBool(ui.IsAbs))
	vm.SetVar("@X", script.NewNumber(float64(ui.X)))
	vm.SetVar("@Y", script.NewNumber(float64(ui.Y)))
	vm.SetVar("@Width", script.NewNumber(float64(ui.Width)))
	vm.SetVar("@Height", script.NewNumber(float64(ui.Height)))
	vm.SetVar("@SelectNo", script.NewNumber(float64(ui.SelectNo)))
	vm.SetVar("@Action", script.NewString(ui.Action))
}

func (ui *UIBase) LoadFromVM(vm *script.VM) {
	ui.Count = int(vm.GetVar("@Count").Num)
	ui.IsEnable = vm.GetVar("@IsEnable").Bool
	ui.IsVisivle = vm.GetVar("@IsVisivle").Bool
	ui.IsAbs = vm.GetVar("@IsAbs").Bool
	ui.X = int(vm.GetVar("@X").Num)
	ui.Y = int(vm.GetVar("@Y").Num)
	ui.Width = int(vm.GetVar("@Width").Num)
	ui.Height = int(vm.GetVar("@Height").Num)
	ui.SelectNo = int(vm.GetVar("@SelectNo").Num)
	ui.Action = vm.GetVar("@Action").Str
}

func (ui *UIBase) Init(type_ UIType) {
	ui.Type = type_
	ui.IsEnable = true
	ui.IsVisivle = true
	// 0だと設定忘れ時どこがおかしいかわからないのでデフォルト値を入れる
	ui.Width = 80
	ui.Height = 24
}

func (ui *UIBase) SetScript(scriptSrc string) error {
	var err error
	ui.script, err = script.Compile(scriptSrc)
	return err
}

func (ui *UIBase) Update() {
	if ui.script != nil {
		ui.StoreToVM(ui.script.GetVM())
		ui.script.Run()
		ui.LoadFromVM(ui.script.GetVM())
	}
}
