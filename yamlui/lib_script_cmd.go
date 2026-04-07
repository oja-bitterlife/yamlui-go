package yamlui

import (
	"path"
	"strings"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// scriptコマンドを登録する
func (lib *YAMLUI) SetUIScriptCmds() {
	// UIBaseのscriptコマンドを登録する
	script.RegisterVMCmd("event", lib.scriptEvent)
	script.RegisterVMCmd("action", lib.scriptAction)
}

// ==================================================
// eventコマンドをscriptに登録する
func (lib *YAMLUI) scriptEvent(vm *script.VM, args []script.Value) (script.Value, error) {
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
	for k := range vm.Vars {
		event, ok := strings.CutPrefix(k, "@")
		if !ok {
			continue // プレフィックスが違うものはスキップ
		}

		// path.Matchでイベント名をマッチさせる
		if match, err := path.Match(value.Str, event); err == nil && match {
			return script.NewBool(true), nil // マッチしたらtrueを返す
		}
	}
	return script.NewBool(false), nil // マッチしなかったらfalseを返す
}

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
		return script.Value{}, script.LogErr("action command expects a string argument, but got " + value.Type.String())
	}

	lib.AddEvent(value.Str) // イベントを追加

	return value, nil
}
