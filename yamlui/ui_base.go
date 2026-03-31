package yamlui

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// UIの基本構造.これを保有して各UI構造体を作る
type UIBase struct {
	// Scriptで更新させないもの
	Type  string
	ID    string
	Frame int

	// スクリプトで更新できるプロパティ
	// ----------------------------------------
	IsEnable bool

	// 座標
	IsAbs bool
	X     float64
	Y     float64
	W     float64
	H     float64

	// 表示
	IsVisible bool
	Text      string
	Color     string // 使用するカラーの名称。system/msg/frameなどを想定

	// インタラクティブなUIに必要なプロパティ
	SelectNo float64
	SelGridX float64 // SelectGridで横の折り返し位置
	Action   string  // 都度リセットされる

	// 子要素が保存させたいもの
	Prop map[string]script.Value

	// 保存しないもの
	// ----------------------------------------
	children []*UIBase
	script   *script.Runtime

	// インターフェース(ランタイム用)
	// ----------------------------------------
	onInitIF     OnInitIF
	updateIF     UpdateIF
	updateTreeIF UpdateTreeIF
	drawIF       DrawIF
	drawTreeIF   DrawTreeIF
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
	vm.SetVar("@IsVisivle", script.NewBool(self.IsVisible))
	vm.SetVar("@Text", script.NewString(self.Text))
	vm.SetVar("@Color", script.NewString(self.Color))
	vm.SetVar("@SelectNo", script.NewNumber(float64(self.SelectNo)))
	vm.SetVar("@SelGridX", script.NewNumber(float64(self.SelGridX)))
	vm.SetVar("@Action", script.NewString(self.Action))
	for k, v := range self.Prop {
		vm.SetVar("@Prop."+k, v)
	}
}

func (self *UIBase) loadFromVM(vm *script.VM) {
	self.IsEnable = vm.GetVar("@IsEnable").Bool
	self.IsAbs = vm.GetVar("@IsAbs").Bool
	self.X = vm.GetVar("@X").Num
	self.Y = vm.GetVar("@Y").Num
	self.W = vm.GetVar("@Width").Num
	self.H = vm.GetVar("@Height").Num
	self.IsVisible = vm.GetVar("@IsVisivle").Bool
	self.Text = vm.GetVar("@Text").Str
	self.Color = vm.GetVar("@Color").Str
	self.SelectNo = vm.GetVar("@SelectNo").Num
	self.SelGridX = vm.GetVar("@SelGridX").Num
	self.Action = vm.GetVar("@Action").Str
	for k, v := range vm.GetVars() {
		prefix := "@Prop."
		if after, ok := strings.CutPrefix(k, prefix); ok {
			propKey := after
			self.Prop[propKey] = v
		}
	}
}

// **********************************************************************
// UIBaseの関数
func NewUIBase(type_ string) *UIBase {
	// 仮のIDを生成
	b := make([]byte, 16) // 128bit
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	ui := &UIBase{}
	ui.Type = type_
	ui.ID = hex.EncodeToString(b)
	ui.IsEnable = true
	ui.IsVisible = true
	ui.Color = "system"
	// とりあえず大きな値を入れておく
	ui.W = 65536
	ui.H = 65536
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

// UIBaseをValueに変換する関数
func (self *UIBase) ToValue() script.Value {
	return script.NewLitMap(map[string]script.Value{
		"Type":      script.NewString(self.Type),
		"ID":        script.NewString(self.ID),
		"IsEnable":  script.NewBool(self.IsEnable),
		"IsAbs":     script.NewBool(self.IsAbs),
		"X":         script.NewNumber(float64(self.X)),
		"Y":         script.NewNumber(float64(self.Y)),
		"W":         script.NewNumber(float64(self.W)),
		"H":         script.NewNumber(float64(self.H)),
		"IsVisible": script.NewBool(self.IsVisible),
		"Text":      script.NewString(self.Text),
		"Color":     script.NewString(self.Color),
		"SelectNo":  script.NewNumber(float64(self.SelectNo)),
		"SelGridX":  script.NewNumber(float64(self.SelGridX)),
		"Action":    script.NewString(self.Action),
		"Prop":      script.NewLitMap(self.Prop),
	})
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
	OnInit(ui *UIBase, ctx UpdateContext) error
}

type UpdateIF interface {
	Update(ui *UIBase, ctx UpdateContext) error
}

type UpdateTreeIF interface {
	UpdateTree(ui *UIBase, ctx UpdateContext) error
}

type UpdateContext struct {
	Parent *UIBase
	frame  int
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

func NewUpdateContext(parent *UIBase, frame int) UpdateContext {
	return UpdateContext{
		Parent: parent,
		frame:  frame,
	}
}

// ----------------------------------------
// UpdateIFの呼び出し
func (self *UIBase) callUpdate(ctx UpdateContext) error {
	var err error

	if self.IsEnable {
		if self.updateIF != nil {
			self.updateIF.Update(self, ctx)
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

func (self *UIBase) callUpdateTree(ctx UpdateContext) error {
	if self.updateTreeIF != nil {
		return self.updateTreeIF.UpdateTree(self, ctx)
	} else {
		return self.updateTree(ctx)
	}
}

// ----------------------------------------
// Update実行
// 呼び出し口
func (self *UIBase) Update(frame int) error {
	ctx := NewUpdateContext(self, frame)

	var lastErr error
	if self.Frame == 0 {
		// 最初のフレームならOnInitを呼び出す
		if self.onInitIF != nil {
			lastErr = self.onInitIF.OnInit(self, ctx)
		}
	} else {
		// それ以降のフレームならUpdateを呼び出す
		if err := self.callUpdate(ctx); err != nil {
			lastErr = err
		}
	}

	// InitもUpdateもTreeは手繰る
	if err := self.callUpdateTree(ctx); err != nil {
		lastErr = err
	}

	self.Frame++

	return lastErr
}

// 再帰実行
func (self *UIBase) updateTree(ctx UpdateContext) error {
	ctx = NewUpdateContext(self, ctx.frame)

	var lastErr error

	// 先に子のUpdateを全部実行する
	for _, child := range self.children {
		if err := child.callUpdate(ctx); err != nil {
			lastErr = err
		}
	}

	// そのあとTreeを手繰る
	for _, child := range self.children {
		if err := child.callUpdateTree(ctx); err != nil {
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
	Draw(ui *UIBase, clip Area, ctx DrawContext)
}

// クリップ操作が必要な時とか
type DrawTreeIF interface {
	DrawTree(ui *UIBase, clip Area)
}

type DrawContext struct {
	Parent     *UIBase
	ParentArea Area
}

func (self *UIBase) SetDrawIF(drawIF DrawIF) {
	self.drawIF = drawIF
}

func (self *UIBase) SetDrawTreeIF(drawTreeIF DrawTreeIF) {
	self.drawTreeIF = drawTreeIF
}

func NewDrawContext(parent *UIBase, parentArea Area) DrawContext {
	return DrawContext{
		Parent:     parent,
		ParentArea: parentArea,
	}
}

// ----------------------------------------
// 描画IFの呼び出し
func (self *UIBase) callDraw(clip Area, ctx DrawContext) {
	if self.IsVisible && self.drawIF != nil {
		self.drawIF.Draw(self, clip, ctx)
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
func (self *UIBase) Draw(screen Area) {
	ctx := NewDrawContext(self, screen)
	area := self.calcDrawArea(screen)

	// 先に自分を描画する
	self.callDraw(area, ctx)

	// 子のDrawを呼び出す
	self.callDrawTree(area)
}

// 再帰実行
func (self *UIBase) drawTree(clip Area) {
	ctx := NewDrawContext(self, clip)

	// 先に子のDrawを全部実行する
	for _, child := range self.children {
		if child.drawIF != nil && child.IsVisible {
			area := child.calcDrawArea(clip)
			child.callDraw(area, ctx)
		}
	}

	// そのあとTreeを手繰る
	for _, child := range self.children {
		if child.IsVisible {
			area := child.calcDrawArea(clip)
			child.callDrawTree(area)
		}
	}
}
