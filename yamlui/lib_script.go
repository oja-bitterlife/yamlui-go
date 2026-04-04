package yamlui

import (
	"path"
	"strings"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
const (
	UIEventPrefix = "@UIEvent." // UIイベントのプロパティ名のプレフィックス
)

// vmのPropにUIのイベントをセットする
func (self *UIBase) storeScriptEvent(events []string) {
	for _, event := range events {
		self.script.SetVar(UIEventPrefix+event, script.NewBool(true))
	}
}

// **********************************************************************
// scriptコマンドを登録する
func (self *UIBase) setUIScriptCmds(vm *script.VM) {
	vm.RegisterCmd("event",
		func(vm *script.VM, args []script.Value) (script.Value, error) {
			return self.scriptEvent(vm, args)
		})
}

// ==================================================
// eventコマンドをscriptに登録する
func (ui *UIBase) scriptEvent(vm *script.VM, args []script.Value) (script.Value, error) {
	// 引数の数をチェック
	if err := script.CheckCmdArgNum("event", 1, args); err != nil {
		return script.Value{}, err
	}
	// 引数を評価
	value, err := vm.Eval(args[0])
	if err != nil {
		return script.Value{}, err
	}

	// 文字列のはずなのでチェック
	if value.Type != script.TypeString {
		return script.Value{}, script.LogErr("event command expects a string argument, but got " + value.Type.String())
	}

	// vm.varsの中にあるかpath.Matchでチェックする
	for k := range vm.GetVars() {
		event, ok := strings.CutPrefix(k, UIEventPrefix)
		if !ok {
			continue // プレフィックスが違うものはスキップ
		}

		// path.Matchでイベント名をマッチさせる
		if match, err := path.Match(value.Str, event); err == nil && match {
			vm.SetVar("@Action", value)      // @Actionにイベントをセットしておく
			return script.NewBool(true), nil // マッチしたらtrueを返す
		}
	}
	return script.NewBool(false), nil // マッチしなかったらfalseを返す
}
