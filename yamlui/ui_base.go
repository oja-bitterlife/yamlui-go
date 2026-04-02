package yamlui

import (
	"crypto/rand"
	"encoding/hex"
	"maps"
	"slices"
	"strings"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// UIの基本構造.これを保有して各UI構造体を作る
type UIBase struct {
	// ==================================================
	// Scriptで更新させないもの
	ID          string
	UpdateCount int // @Frameとして送り込む

	// Event
	Events []string

	// ==================================================
	// スクリプトで更新できるプロパティ
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
	// Action   string  // 直接EventをAddすればいいのでActionは廃止

	// 都度リセットされる
	ScriptAction string       // スクリプト中に@Actionへ発生したイベント名を入れる
	ScriptResult script.Value // スクリプトの評価結果

	// 子要素が保存させたいもの
	Prop map[string]script.Value

	// ==================================================
	// 保存しないもの
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

func NewUIBase() *UIBase {
	// 仮のIDを生成
	b := make([]byte, 16) // 128bit
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	ui := &UIBase{}
	ui.ID = hex.EncodeToString(b)
	ui.IsEnable = true
	ui.IsVisible = true
	ui.Color = "system"
	// とりあえず大きな値を入れておく
	ui.W = 65536
	ui.H = 65536

	// map/sliceを初期化しておく
	ui.Events = []string{}
	ui.Prop = make(map[string]script.Value)
	ui.children = []*UIBase{}

	return ui
}

// **********************************************************************
// UIComponentIFの実装
func (self *UIBase) GetUIBase() *UIBase {
	return self
}

func (self *UIBase) Clone() *UIBase {
	clone := *self

	// PropとEventsは標準ライブラリのCloneでOK
	clone.Prop = maps.Clone(self.Prop)
	clone.Events = slices.Clone(self.Events)

	// UIBaseのポインタは再帰的にUICloneする
	if self.children != nil {
		clone.children = make([]*UIBase, len(self.children))
		for i, child := range self.children {
			// 子要素の UIClone を呼び、新しい個体として登録する
			clone.children[i] = child.Clone()
		}
	}

	// scriptはRuntimeがCloneを持っている
	if self.script != nil {
		clone.script = self.script.Clone()
	}

	return &clone
}

func (self *UIBase) Setup(lib *YAMLUI, type_ string, parent *UIBase, data map[string]script.Value) error {
	// 基本的にはUIBaseは何もしない。今後のための予約
	return nil
}

// **********************************************************************
// VMとのやりとり
func (self *UIBase) storeToVM(vm *script.VM) {
	// スクリプトで使うFrameはUpdateCountを入れる
	vm.SetVar("@Frame", script.NewNumber(float64(self.UpdateCount)))

	// プロパティの送信
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

	// スクリプトからUIに伝えるためのものなので常に空文字を入れておく
	vm.SetVar("@Action", script.NewString(""))

	// PropはProp.をプレフィックスにして送る
	for k, v := range self.Prop {
		vm.SetVar("@Prop."+k, v)
	}
}

func (self *UIBase) loadFromVM(vm *script.VM) {
	// プロパティの受信
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
	self.ScriptAction = vm.GetVar("@Action").Str

	// PropはProp.をプレフィックスにして受け取る
	for k, v := range vm.GetVars() {
		if after, ok := strings.CutPrefix(k, "@Prop."); ok {
			propKey := after
			self.Prop[propKey] = v
		}
	}
}

// **********************************************************************
// UIBaseの関数
func (self *UIBase) SetScript(scriptSrc string) error {
	var err error
	self.script, err = script.Compile(scriptSrc)
	return err
}

func (self *UIBase) GetScriptRuntime() *script.Runtime {
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

func (self *UIBase) FindChildByID(id string) *UIBase {
	for _, c := range self.children {
		if c.ID == id {
			return c
		}
		if found := c.FindChildByID(id); found != nil {
			return found
		}
	}
	return nil
}
