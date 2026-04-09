package yamlui

import (
	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// scriptコマンドを登録する
func (lib *YAMLUI) SetUIScriptCmds() {
	// UIBaseのscriptコマンドを登録する
	lib.RegisterVMCmd("action", lib.scriptAction)
}

// ==================================================
// eventコマンドをscriptに登録する
func (lib *YAMLUI) scriptAction(vm *script.VM, args []script.Value) (script.Value, error) {
	// 引数の数をチェック
	if err := script.CheckCmdArgNum("action", 1, args); err != nil {
		return script.Value{}, err
	}

	// 引数を評価
	value, err := vm.Eval(args[0])
	if err != nil {
		return script.Value{}, err
	}

	// 文字列のはずなのでチェック
	if value.Type != script.TypeString {
		return script.Value{}, lib.LogErr("action command expects a string argument, but got " + value.Type.String())
	}

	lib.SendEvent(value.Str) // イベントを追加

	return value, nil
}
