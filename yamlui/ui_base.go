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

	onInitIF     OnInitIF
	updateIF     UpdateIF
	updateTreeIF UpdateTreeIF
	drawIF       DrawIF
	drawTreeIF   DrawTreeIF

	script *script.Runtime
}

// ==================================================
// VMとのやりとり
func (self *UIBase) storeToVM(vm *script.VM) {
	vm.SetVar("@Frame", script.NewNumber(float64(self.Frame)))
	vm.SetVar("@IsEnable", script.NewBool(self.IsEnable))
	vm.SetVar("@IsAbs", script.NewBool(self.IsAbs))
	vm.SetVar("@X", script.NewNumber(float64(self.X)))
	vm.SetVar("@Y", script.NewNumber(float64(self.Y)))
	vm.SetVar("@Width", script.NewNumber(float64(self.W)))
	vm.SetVar("@Height", script.NewNumber(float64(self.H)))
	vm.SetVar("@IsVisivle", script.NewBool(self.IsVisivle))
	vm.SetVar("@Text", script.NewString(self.Text))
	vm.SetVar("@Color", script.NewString(self.Color))
	vm.SetVar("@SelectNo", script.NewNumber(float64(self.SelectNo)))
	vm.SetVar("@SelGridX", script.NewNumber(float64(self.SelGridX)))
	// 都度リセット
	vm.SetVar("@Action", script.NewString(""))
}

func (self *UIBase) loadFromVM(vm *script.VM) {
	self.IsEnable = vm.GetVar("@IsEnable").Bool
	self.IsAbs = vm.GetVar("@IsAbs").Bool
	self.X = int(vm.GetVar("@X").Num)
	self.Y = int(vm.GetVar("@Y").Num)
	self.W = int(vm.GetVar("@Width").Num)
	self.H = int(vm.GetVar("@Height").Num)
	self.IsVisivle = vm.GetVar("@IsVisivle").Bool
	self.Text = vm.GetVar("@Text").Str
	self.Color = vm.GetVar("@Color").Str
	self.SelectNo = int(vm.GetVar("@SelectNo").Num)
	self.SelGridX = int(vm.GetVar("@SelGridX").Num)
	self.Action = vm.GetVar("@Action").Str
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

func (self *UIBase) SetScript(scriptSrc string) error {
	var err error
	self.script, err = script.Compile(scriptSrc)
	return err
}

func (self *UIBase) GetRuntime() *script.Runtime {
	return self.script
}

func (self *UIBase) SetVar(name string, value script.Value) {
	self.script.GetVM().SetVar(name, value)
}

// **********************************************************************
// Tree構造化
// ==================================================
// ツリー操作
func (self *UIBase) AddChild(child *UIBase) {
	self.children = append(self.children, child)
}

func (self *UIBase) RemoveChild(child *UIBase) {
	for i, c := range self.children {
		if c == child {
			self.children = append(self.children[:i], self.children[i+1:]...)
			return
		}
	}
}

// ==================================================
// 更新
// ----------------------------------------
// Updateのインターフェース
type OnInitIF interface {
	OnInit(ui *UIBase) error
}

type UpdateIF interface {
	Update(ui *UIBase, frame int) error
}

type UpdateTreeIF interface {
	UpdateTree(ui *UIBase, frame int) error
}

func (self *UIBase) SetOnInitIF(onInitIF OnInitIF) {
	self.onInitIF = onInitIF
}

func (self *UIBase) SetUpdateIF(updateIF UpdateIF) {
	self.updateIF = updateIF
}

func (self *UIBase) SetUpdateTreeIF(updateTreeIF UpdateTreeIF) {
	self.updateTreeIF = updateTreeIF
}

// ----------------------------------------
// UpdateIFの呼び出し
func (self *UIBase) callUpdate(frame int) error {
	var err error

	if self.IsEnable {
		if self.updateIF != nil {
			self.updateIF.Update(self, frame)
		}

		// UIの更新後スクリプトがあれば走らせる
		if self.script != nil {
			self.storeToVM(self.script.GetVM())
			_, err = self.script.Run()
			self.loadFromVM(self.script.GetVM())
		}
	}

	return err
}

func (self *UIBase) callUpdateTree(frame int) error {
	if self.updateTreeIF != nil {
		return self.updateTreeIF.UpdateTree(self, frame)
	} else {
		return self.updateTree(frame)
	}
}

// ----------------------------------------
// Update実行
// 呼び出し口
func (self *UIBase) Update(frame int) error {
	var lastErr error
	if self.Frame == 0 {
		// 最初のフレームならOnInitを呼び出す
		if self.onInitIF != nil {
			lastErr = self.onInitIF.OnInit(self)
		}
	} else {
		// それ以降のフレームならUpdateを呼び出す
		if err := self.callUpdate(frame); err != nil {
			lastErr = err
		}
	}

	// InitもUpdateもTreeは手繰る
	if err := self.callUpdateTree(frame); err != nil {
		lastErr = err
	}

	self.Frame++

	return lastErr
}

// 再帰実行
func (self *UIBase) updateTree(frame int) error {
	var lastErr error

	// 先に子のUpdateを全部実行する
	for _, child := range self.children {
		if err := child.callUpdate(frame); err != nil {
			lastErr = err
		}
	}

	// そのあとTreeを手繰る
	for _, child := range self.children {
		if err := child.callUpdateTree(frame); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// ==================================================
// 描画
// ----------------------------------------
// 描画のインターフェース
// 描画を行うとき
type DrawIF interface {
	Draw(ui *UIBase, clip Area)
}

// クリップ操作が必要な時とか
type DrawTreeIF interface {
	DrawTree(ui *UIBase, clip Area)
}

func (self *UIBase) SetDrawIF(drawIF DrawIF) {
	self.drawIF = drawIF
}

func (self *UIBase) SetDrawTreeIF(drawTreeIF DrawTreeIF) {
	self.drawTreeIF = drawTreeIF
}

// ----------------------------------------
// 描画IFの呼び出し
func (self *UIBase) callDraw(clip Area) {
	if self.IsVisivle && self.drawIF != nil {
		self.drawIF.Draw(self, clip)
	}
}

func (self *UIBase) callDrawTree(clip Area) {
	if self.drawTreeIF != nil {
		self.drawTreeIF.DrawTree(self, clip)
	} else {
		self.drawTree(clip)
	}
}

// ----------------------------------------
// 描画実行
func (self *UIBase) calcDrawArea(clip Area) Area {
	if self.IsAbs {
		return Area{
			X: self.X,
			Y: self.Y,
			W: self.W,
			H: self.H,
		}.Clip(clip)
	} else {
		return Area{
			X: clip.X + self.X,
			Y: clip.Y + self.Y,
			W: self.W,
			H: self.H,
		}.Clip(clip)
	}
}

// 呼び出し口
func (self *UIBase) Draw(clip Area) {
	area := self.calcDrawArea(clip)

	// 先に自分を描画する
	self.callDraw(area)

	// 子のDrawを呼び出す
	self.callDrawTree(area)
}

// 再帰実行
func (self *UIBase) drawTree(clip Area) {
	// 先に子のDrawを全部実行する
	for _, child := range self.children {
		if child.drawIF != nil && child.IsVisivle {
			area := child.calcDrawArea(clip)
			child.callDraw(area)
		}
	}

	// そのあとTreeを手繰る
	for _, child := range self.children {
		if child.IsVisivle {
			area := child.calcDrawArea(clip)
			child.callDrawTree(area)
		}
	}
}
