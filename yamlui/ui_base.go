package yamlui

import (
	"maps"
	"slices"
	"strings"
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

	// 子要素が保存させたいもの
	Prop map[string]script.Value // PropはValueMap型と意味が違うのでmapそのままで

	// ==================================================
	// 保存しないもの
	children []*UIBase
	script   *script.VM

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
	lastID++
	ui.IsEnable = true
	ui.IsVisible = true
	// とりあえず大きな値を入れておく
	ui.W = 65536
	ui.H = 65536

	// map/sliceを初期化しておく
	ui.Events = []string{}
	ui.Prop = make(map[string]script.Value) // ValueMapにはしない
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

func (self *UIBase) Setup(type_ string, data script.ValueMap) error {
	// 基本的にはUIBaseは何もしない。今後のための予約
	return nil
}

// **********************************************************************
// ScriptVMとの連携
func (self *UIBase) storeToVM(vm *script.VM) {
	// 既存の変数をクリア
	vm.ClearVars()

	// スクリプトで使うFrameはUpdateCountを入れる
	vm.SetVar("@Frame", script.NewNumber(self.UpdateCount))

	// プロパティの送信
	vm.SetVar("@IsEnable", script.NewBool(self.IsEnable))
	vm.SetVar("@Remove", script.NewBool(self.Remove))
	vm.SetVar("@X", script.NewNumber(self.X))
	vm.SetVar("@Y", script.NewNumber(self.Y))
	vm.SetVar("@Width", script.NewNumber(self.W))
	vm.SetVar("@Height", script.NewNumber(self.H))
	vm.SetVar("@IsVisivle", script.NewBool(self.IsVisible))
	vm.SetVar("@Text", script.NewString(self.Text))

	// Propを送る
	for k, v := range self.Prop {
		vm.SetVar("@"+k, v)
	}
}

func (self *UIBase) loadFromVM(vm *script.VM) {
	// プロパティの受信
	for k, v := range vm.GetVars() {
		// @で始まる変数はUIBaseのプロパティとして受け取る
		if propName, ok := strings.CutPrefix(k, "@"); ok {
			switch propName {

			// UIBaseのプロパティの受信
			case "IsEnable":
				self.IsEnable = v.Bool
			case "Remove":
				self.Remove = v.Bool
			case "X":
				self.X = v.Num
			case "Y":
				self.Y = v.Num
			case "Width":
				self.W = v.Num
			case "Height":
				self.H = v.Num
			case "IsVisivle":
				self.IsVisible = v.Bool
			case "Text":
				self.Text = v.Str

			// その他はProp
			default:
				// "UIEvent."で始まるプロパティはイベントなのでPropには入れない
				if !strings.HasPrefix(propName, UIEventPrefix) {
					// 一般的なProperty
					self.Prop[propName] = v
				}
			}
		}
	}
}

// ==================================================
// getter/setter
func (self *UIBase) SetScriptVM(newVM *script.VM) {
	self.script = newVM
}

func (self *UIBase) GetScriptVM() (*script.VM, error) {
	if self.script == nil {
		return nil, script.LogErr("No script runtime set for this UIBase (ID: " + self.ID + ")")
	}
	return self.script, nil
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
