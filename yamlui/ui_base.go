package yamlui

import "github.com/oja-bitterlife/yamlui-go/script"

// **********************************************************************
// UIの基本構造.これを保有して各UI構造体を作る
type UIBase struct {
	Type  string
	Frame int

	// スクリプトで更新できるプロパティ
	// ----------------------------------------
	IsEnable bool

	// 座標
	IsAbs  bool
	X      int
	Y      int
	Width  int
	Height int

	// 表示
	IsVisivle bool
	Text      string
	Color     string

	// インタラクティブなUIに必要なプロパティ
	SelectNo int
	SelGridX int // SelectGridで横の折り返し位置
	Action   string

	// 保存しないもの
	// ----------------------------------------
	script     *script.Runtime
	onInitFunc func(ui *UIBase, x, y int)
	updateFunc func(ui *UIBase, x, y int)
	drawFunc   func(ui *UIBase, x, y int)
}

// ==================================================
// VMとのやりとり
func (ui *UIBase) storeToVM(vm *script.VM) {
	vm.SetVar("@Frame", script.NewNumber(float64(ui.Frame)))
	vm.SetVar("@IsEnable", script.NewBool(ui.IsEnable))
	vm.SetVar("@IsAbs", script.NewBool(ui.IsAbs))
	vm.SetVar("@X", script.NewNumber(float64(ui.X)))
	vm.SetVar("@Y", script.NewNumber(float64(ui.Y)))
	vm.SetVar("@Width", script.NewNumber(float64(ui.Width)))
	vm.SetVar("@Height", script.NewNumber(float64(ui.Height)))
	vm.SetVar("@IsVisivle", script.NewBool(ui.IsVisivle))
	vm.SetVar("@Text", script.NewString(ui.Text))
	vm.SetVar("@Color", script.NewString(ui.Color))
	vm.SetVar("@SelectNo", script.NewNumber(float64(ui.SelectNo)))
	vm.SetVar("@SelGridX", script.NewNumber(float64(ui.SelGridX)))
	vm.SetVar("@Action", script.NewString(ui.Action))
}

func (ui *UIBase) loadFromVM(vm *script.VM) {
	ui.Frame = int(vm.GetVar("@Frame").Num)
	ui.IsEnable = vm.GetVar("@IsEnable").Bool
	ui.IsAbs = vm.GetVar("@IsAbs").Bool
	ui.X = int(vm.GetVar("@X").Num)
	ui.Y = int(vm.GetVar("@Y").Num)
	ui.Width = int(vm.GetVar("@Width").Num)
	ui.Height = int(vm.GetVar("@Height").Num)
	ui.IsVisivle = vm.GetVar("@IsVisivle").Bool
	ui.Text = vm.GetVar("@Text").Str
	ui.Color = vm.GetVar("@Color").Str
	ui.SelectNo = int(vm.GetVar("@SelectNo").Num)
	ui.SelGridX = int(vm.GetVar("@SelGridX").Num)
	ui.Action = vm.GetVar("@Action").Str
}

// **********************************************************************
// UIBaseの関数
func NewUIBase(type_ string) *UIBase {
	ui := &UIBase{}
	ui.Type = type_
	ui.IsEnable = true
	ui.IsVisivle = true
	ui.Color = "system"
	// 0だと設定忘れ時どこがおかしいかわからないのでデフォルト値を入れる
	ui.Width = 80
	ui.Height = 24
	return ui
}

func (ui *UIBase) SetScript(scriptSrc string) error {
	var err error
	ui.script, err = script.Compile(scriptSrc)
	return err
}

func (ui *UIBase) SetOnInitFunc(initFunc func(ui *UIBase, x, y int)) {
	ui.onInitFunc = initFunc
}

func (ui *UIBase) SetUpdateFunc(updateFunc func(ui *UIBase, x, y int)) {
	ui.updateFunc = updateFunc
}

func (ui *UIBase) SetDrawFunc(drawFunc func(ui *UIBase, x, y int)) {
	ui.drawFunc = drawFunc
}

// **********************************************************************
// 実行
func (ui *UIBase) Update() {
	if ui.IsEnable {
		if ui.Frame == 0 && ui.onInitFunc != nil {
			ui.onInitFunc(ui, ui.X, ui.Y)
		}
		ui.Frame++

		if ui.updateFunc != nil {
			ui.updateFunc(ui, ui.X, ui.Y)
		}

		if ui.script != nil {
			ui.storeToVM(ui.script.GetVM())
			ui.script.Run()
			ui.loadFromVM(ui.script.GetVM())
		}
	}

	if ui.IsVisivle && ui.drawFunc != nil {
		ui.drawFunc(ui, ui.X, ui.Y)
	}
}
