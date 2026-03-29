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

	onInitIF Initarizable
	updateIF Updatable
	drawIF   Drawable

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

type Initarizable interface {
	OnInit(ui *UIBase)
}

func (ui *UIBase) SetOnInitFunc(onInitIF Initarizable) {
	ui.onInitIF = onInitIF
}

func (ui *UIBase) SetUpdateIF(updateIF Updatable) {
	ui.updateIF = updateIF
}

func (ui *UIBase) SetDrawIF(drawIF Drawable) {
	ui.drawIF = drawIF
}

func (ui *UIBase) GetRuntime() *script.Runtime {
	return ui.script
}

// **********************************************************************
// Tree構造化
// ==================================================
// ツリー操作
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

// ==================================================
// 更新
type Updatable interface {
	Update(ui *UIBase, frame int) error
}

func (ui *UIBase) update(frame int) error {
	var err error

	if ui.IsEnable {
		if ui.Frame == 0 && ui.onInitIF != nil {
			ui.onInitIF.OnInit(ui)
		}
		ui.Frame++

		if ui.updateIF != nil {
			ui.updateIF.Update(ui, frame)
		}

		if ui.script != nil {
			ui.storeToVM(ui.script.GetVM())
			_, err = ui.script.Run()
			ui.loadFromVM(ui.script.GetVM())
		}
	}

	return err
}

// 呼び出し口
func (ui *UIBase) UpdateTree(frame int) error {
	lastErr := ui.update(frame)

	if err := ui.updateTree(frame); err != nil {
		lastErr = err
	}

	return lastErr
}

// 再帰実行
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

// ==================================================
// 描画
type Drawable interface {
	Draw(ui *UIBase, x, y int)
}

// 呼び出し口
func (ui *UIBase) DrawTree(x, y int) {
	if ui.drawIF != nil {
		ui.drawIF.Draw(ui, x, y)
	}

	ui.drawTree(x, y)
}

// 再帰実行
func (ui *UIBase) drawTree(x, y int) {
	// 先に子のDrawを全部実行する
	for _, child := range ui.children {
		if child.drawIF != nil {
			child.drawIF.Draw(ui, x, y)
		}
	}

	// そのあとTreeを手繰る
	for _, child := range ui.children {
		child.drawTree(x, y)
	}
}
