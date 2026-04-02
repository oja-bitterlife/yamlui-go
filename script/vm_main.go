package script

// **********************************************************************
// 評価器(VM)の実装
// ==================================================
// 外部との連携用データ構造体
type VM struct {
	vars         Value
	cmds         map[string]func(vm *VM, args []Value) (Value, error)
	maxRecursion int // 再帰の最大深さ
	maxRepeat    int // repeatの最大回数
}

// VMの初期化
func NewVM() *VM {
	return &VM{
		vars:         NewLitMap(make(map[string]Value)),
		cmds:         make(map[string]func(vm *VM, args []Value) (Value, error)),
		maxRecursion: 64,
		maxRepeat:    256,
	}
}

// ==================================================
// コマンドを登録する関数
func (vm *VM) RegisterCmd(name string, fn func(vm *VM, args []Value) (Value, error)) {
	vm.cmds[name] = fn
}

// mapを渡して複数のコマンドを登録する
func (vm *VM) RegisterCmdList(cmds map[string]func(vm *VM, args []Value) (Value, error)) {
	for name, fn := range cmds {
		vm.RegisterCmd(name, fn)
	}
}

// デバッグ用cmds取得関数
func (vm *VM) GetCmds() map[string]func(vm *VM, args []Value) (Value, error) {
	return vm.cmds
}

// ==================================================
// vars用の取得・設定する関数
func (vm *VM) GetVar(name string) Value {
	return vm.vars.Map[name]
}

func (vm *VM) GetVars() map[string]Value {
	return vm.vars.Map
}

func (vm *VM) SetVar(name string, value Value) {
	vm.vars.Map[name] = value
}
