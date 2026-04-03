package yamlui

import "github.com/oja-bitterlife/yamlui-go/script"

func (self *UIBase) setUIScriptCmds() {
	self.script.RegisterCmd("event", scriptEvent)
}

// ==================================================
// eventコマンドをscriptに登録する
func scriptEvent(vm *script.VM, args []script.Value) (script.Value, error) {
	// 引数の数をチェック
	if err := script.CheckCmdArgNum("event", 1, args); err != nil {
		return script.Value{}, err
	}
	// 引数を評価
	value, err := vm.Eval(args[0])
	if err != nil {
		return script.Value{}, err
	}

	// イベントを追加
	vm.SetVar("@Action", value) // @Actionにイベントをセットしておく

	return value, nil
}
