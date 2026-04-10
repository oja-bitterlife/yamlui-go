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
	ID string

	// ScriptVM. VarsはUIBaseのプロパティと共有
	script script.VM

	// Event
	Events []string // キャッチするイベント

	// 保存しないもの
	// ==================================================
	children []*UIBase

	// インターフェース(ランタイム用)
	// ----------------------------------------
	dispatchIF DispatchIF
	drawIF     DrawIF
	drawTreeIF DrawContextIF // z順を変更したいときに使う
}

// IDが決まっているときはこっち
func NewUIBaseWithID(lib *YAMLUI, id string) *UIBase {
	ui := &UIBase{}
	ui.ID = id

	// ScriptVM.cmdsIDはルートノードのIDを入れる。
	// ルートノードがないときはidがRootのIDなのでそれを使う
	var cmdsID string
	if lib.root != nil {
		cmdsID = lib.root.ID
	} else {
		if !script.HasPrefix(id, ROOT_NODE_PREFIX) {
			lib.LogErr("ID does not start with %s: %s", ROOT_NODE_PREFIX, id)
		}
		cmdsID = id
	}
	ui.script = script.NewVM(cmdsID)

	// map/sliceを初期化しておく
	ui.Events = []string{}
	ui.children = []*UIBase{}

	// 0以外の初期値はVarsに入れておく
	// ----------------------------------------
	ui.SetPropBool(PROP_IS_ENABLE, true)
	ui.SetPropBool(PROP_IS_VISIBLE, true)
	ui.SetPropNum(PROP_W, 65536)
	ui.SetPropNum(PROP_H, 65536)

	return ui
}

// 匿名板。だいたいこっちでOK
func NewUIBase(lib *YAMLUI) *UIBase {
	newID := atomic.AddInt32(&lastID, 1)
	return NewUIBaseWithID(lib, "UIBase_"+script.Itoa(int(newID)))
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

// **********************************************************************
// ScriptVMとの連携
// 参照用で送る
func (self *UIBase) storeToVM(frame int, event string) {
	self.script.Vars[PROP_ID] = script.NewString(self.ID)
	self.script.Vars[PROP_EVENT] = script.NewString(event) // 発生したイベント
	self.script.Vars[PROP_FRAME] = script.NewNumber(frame)
}

// ==================================================
// getter/setter
func (self *UIBase) HasScript() bool {
	return len(self.script.AST.List) > 0
}

// VM単体でテストしたいときにアクセスするための関数
func (self *UIBase) GetVM() *script.VM {
	return &self.script
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
