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
	IsAbs bool
	X     int
	Y     int
	W     int
	H     int

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
	vm.SetVar("@Width", script.NewNumber(float64(ui.W)))
	vm.SetVar("@Height", script.NewNumber(float64(ui.H)))
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
	ui.W = int(vm.GetVar("@Width").Num)
	ui.H = int(vm.GetVar("@Height").Num)
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
	ui.W = 64
	ui.H = 48 / 2
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
		err := child.updateTree(frame)
		if err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// ==================================================
// 描画
type Drawable interface {
	Draw(ui *UIBase, clip Area)
}

func (ui *UIBase) calcDrawArea(clip Area) Area {
	if ui.IsAbs {
		return Area{
			X: ui.X,
			Y: ui.Y,
			W: ui.W,
			H: ui.H,
		}.Clip(clip)
	} else {
		return Area{
			X: clip.X + ui.X,
			Y: clip.Y + ui.Y,
			W: ui.W,
			H: ui.H,
		}.Clip(clip)
	}
}

// 呼び出し口
func (ui *UIBase) DrawTree(clip Area) {
	if !ui.IsVisivle {
		return
	}

	// 自身の位置で更新
	area := ui.calcDrawArea(clip)

	// 自分のDrawを呼び出す
	if ui.drawIF != nil {
		ui.drawIF.Draw(ui, area)
	}

	// 子のDrawを呼び出す
	ui.drawTree(area)
}

// 再帰実行
func (ui *UIBase) drawTree(clip Area) {
	// 先に子のDrawを全部実行する
	for _, child := range ui.children {
		if child.drawIF != nil && child.IsVisivle {
			area := child.calcDrawArea(clip)
			child.drawIF.Draw(child, area)
		}
	}

	// そのあとTreeを手繰る
	for _, child := range ui.children {
		if child.IsVisivle {
			area := child.calcDrawArea(clip)
			child.drawTree(area)
		}
	}
}
