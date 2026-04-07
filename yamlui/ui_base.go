package yamlui

import (
	"slices"
	"sync/atomic"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// UIの基本構造.これを保有して各UI構造体を作る
var lastID int32

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
	Remove   bool // UIツリーから削除するかどうか。UPDATEの最後にまとめて削除するためのフラグ

	// 座標
	X int
	Y int
	W int
	H int

	// 表示
	IsVisible bool
	Text      string

	// ScriptVM. VarsはUIBaseのプロパティと共有
	script script.VM

	// ==================================================
	// 保存しないもの
	children []*UIBase

	// インターフェース(ランタイム用)
	// ----------------------------------------
	onInitIF     OnInitIF
	updateIF     UpdateIF
	updateTreeIF UpdateTreeIF
	drawIF       DrawIF
	drawTreeIF   DrawTreeIF
}

func NewUIBase() *UIBase {
	ui := &UIBase{}
	newID := atomic.AddInt32(&lastID, 1)
	ui.ID = "UIBase_" + script.Itoa(int(newID))
	ui.IsEnable = true
	ui.IsVisible = true
	// とりあえず大きな値を入れておく
	ui.W = 65536
	ui.H = 65536

	// ScriptVM
	ui.script = script.NewVM(script.Value{})

	// map/sliceを初期化しておく
	ui.Events = []string{}
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
	clone.Events = slices.Clone(self.Events)

	// UIBaseのポインタは再帰的にUICloneする
	if self.children != nil {
		clone.children = make([]*UIBase, len(self.children))
		for i, child := range self.children {
			// 子要素の UIClone を呼び、新しい個体として登録する
			clone.children[i] = child.Clone()
		}
	}

	// scriptはVMがCloneを持っている
	clone.script = self.script.Clone()

	return &clone
}

func (self *UIBase) Setup(type_ string, data script.ValueMap) error {
	// 基本的にはUIBaseは何もしない。今後のための予約
	return nil
}

// **********************************************************************
// ScriptVMとの連携
func (self *UIBase) storeToVM() {
	// スクリプトで使うFrameはUpdateCountを入れる
	self.script.Vars["@Frame"] = script.NewNumber(self.UpdateCount)

	// プロパティの送信
	self.script.Vars["@IsEnable"] = script.NewBool(self.IsEnable)
	self.script.Vars["@Remove"] = script.NewBool(self.Remove)
	self.script.Vars["@X"] = script.NewNumber(self.X)
	self.script.Vars["@Y"] = script.NewNumber(self.Y)
	self.script.Vars["@Width"] = script.NewNumber(self.W)
	self.script.Vars["@Height"] = script.NewNumber(self.H)
	self.script.Vars["@IsVisivle"] = script.NewBool(self.IsVisible)
	self.script.Vars["@Text"] = script.NewString(self.Text)
}

func (self *UIBase) loadFromVM() {
	// UIBaseのプロパティの受信
	for k, v := range self.script.Vars {
		switch k {

		case "@IsEnable":
			self.IsEnable = v.Bool()
		case "@Remove":
			self.Remove = v.Bool()
		case "@X":
			self.X = v.Num
		case "@Y":
			self.Y = v.Num
		case "@Width":
			self.W = v.Num
		case "@Height":
			self.H = v.Num
		case "@IsVisivle":
			self.IsVisible = v.Bool()
		case "@Text":
			self.Text = v.Str
		}
	}
}

// ==================================================
// getter/setter
func (self *UIBase) HasScript() bool {
	return len(self.script.AST.List) > 0
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
