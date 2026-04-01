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
	X float64
	Y float64
	W float64
	H float64

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
