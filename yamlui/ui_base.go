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
	SelGridX int    // SelectGridで横の折り返し位置
	Action   string // 都度リセットされる

	// 保存しないもの
	// ----------------------------------------
	children []*UIBase

	onInitFunc func(ui *UIBase)
	updateFunc func(ui *UIBase, frame int)
	drawFunc   func(ui *UIBase, x, y int)

	script *script.Runtime
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
	// 都度リセット
	vm.SetVar("@Action", script.NewString(""))
}

func (ui *UIBase) loadFromVM(vm *script.VM) {
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
	ui.Width = 64
	ui.Height = 48 / 2
	return ui
}

func (ui *UIBase) SetScript(scriptSrc string) error {
	var err error
	ui.script, err = script.Compile(scriptSrc)
	return err
}

func (ui *UIBase) SetOnInitFunc(initFunc func(*UIBase)) {
	ui.onInitFunc = initFunc
}

func (ui *UIBase) SetUpdateFunc(updateFunc func(*UIBase, int)) {
	ui.updateFunc = updateFunc
}

func (ui *UIBase) SetDrawFunc(drawFunc func(*UIBase, int, int)) {
	ui.drawFunc = drawFunc
}

func (ui *UIBase) GetRuntime() *script.Runtime {
	return ui.script
}

// **********************************************************************
// 実行
func (ui *UIBase) update(frame int) error {
	var err error

	if ui.IsEnable {
		if ui.Frame == 0 && ui.onInitFunc != nil {
			ui.onInitFunc(ui)
		}
		ui.Frame++

		if ui.updateFunc != nil {
			ui.updateFunc(ui, frame)
		}

		if ui.script != nil {
			ui.storeToVM(ui.script.GetVM())
			_, err = ui.script.Run()
			ui.loadFromVM(ui.script.GetVM())
		}
	}

	if ui.IsVisivle && ui.drawFunc != nil {
		ui.drawFunc(ui, ui.X, ui.Y)
	}

	return err
}

func (ui *UIBase) draw(x, y int) {
	if ui.IsVisivle && ui.drawFunc != nil {
		ui.drawFunc(ui, x, y)
	}
}

// **********************************************************************
// Tree構造化
func (ui *UIBase) AddChild(child *UIBase) {
	ui.children = append(ui.children, child)
}

func (ui *UIBase) RemoveChild(child *UIBase) {
	for i, c := range ui.children {
		if c == child {
			ui.children = append(ui.children[:i], ui.children[i+1:]...)
			return
		}
	}
}

func (ui *UIBase) UpdateTree(frame int) error {
	lastErr := ui.update(frame)

	if err := ui.updateTree(frame); err != nil {
		lastErr = err
	}

	return lastErr
}
func (ui *UIBase) updateTree(frame int) error {
	var lastErr error

	// 先に子のUpdateを全部実行する
	for _, child := range ui.children {
		err := child.update(frame)
		if err != nil {
			lastErr = err
		}
	}

	// そのあとTreeを手繰る
	for _, child := range ui.children {
		err := child.UpdateTree(frame)
		if err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (ui *UIBase) DrawTree(x, y int) {
	ui.draw(x, y)
	ui.drawTree(x, y)
}

func (ui *UIBase) drawTree(x, y int) {
	// 先に子のDrawを全部実行する
	for _, child := range ui.children {
		child.draw(x, y)
	}

	// そのあとTreeを手繰る
	for _, child := range ui.children {
		child.drawTree(x, y)
	}
}
