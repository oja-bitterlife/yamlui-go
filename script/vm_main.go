package script

import "sync"

// **********************************************************************
// 評価器(VM)の実装
// ==================================================
const (
	VM_VAR_CAPACITY  = 32 // varsの初期容量
	VM_RECURSION_MAX = 16 // 再帰の最大深さ
	VM_REPEAT_MAX    = 16 // repeatの最大回数
)

// 外部との連携用データ構造体
type VM struct {
	ID     string   // ID
	AST    Value    // コンパイル済みコード
	Vars   ValueMap // VMのメインメモリ的な
	Result Value    // 最後の評価結果を保存する場所
}

// VMの初期化
func NewVM(id string, valueAST Value) VM {
	// VMの初期化
	vm := VM{
		ID:     id,
		AST:    valueAST,
		Vars:   make(map[string]Value, VM_VAR_CAPACITY),
		Result: Value{},
	}

	InitVMCmds(id) // コマンドの初期化

	return vm
}

func (vm *VM) Clone() VM {
	// vars以外はそのままコピーしても問題ない
	clone := *vm

	// varsはCloneする
	for k, v := range vm.Vars {
		clone.Vars[k] = v.Clone()
	}

	return clone
}

// **********************************************************************
// コマンド管理
type VMCmdFunc func(vm *VM, args []Value) (Value, error)

type VMCmds struct {
	systemCmds map[string]VMCmdFunc
	userCmds   map[string]VMCmdFunc
}

var (
	cmds   map[string]VMCmds = make(map[string]VMCmds) // VM IDごとのコマンドセット
	userMu sync.RWMutex
)

func InitVMCmds(id string) {
	userMu.Lock()
	defer userMu.Unlock()

	// すでに初期化されている場合は上書き警告を出す
	if _, exists := cmds[id]; exists {
		LogWarn("VM commands for ID %s are already initialized, overwriting", id)
	}

	// コマンドセットの初期化
	cmds[id] = VMCmds{
		systemCmds: make(map[string]VMCmdFunc),
		userCmds:   make(map[string]VMCmdFunc),
	}
	SetBuiltinCmds(cmds[id].systemCmds) // Builtinの登録
}

func (vm *VM) RegisterVMCmd(cmdName string, fn VMCmdFunc) {
	// 何もしない
	if cmdName == "" || fn == nil {
		LogWarn("invalid command name or function for command registration: name=%s, fn=%v", cmdName, fn)
		return
	}

	// ユーザー定義の登録
	// ----------------------------------------
	userMu.Lock()
	defer userMu.Unlock()

	if _, exists := cmds[vm.ID].userCmds[cmdName]; exists {
		LogWarn("command %s is already defined", cmdName)
	}
	cmds[vm.ID].userCmds[cmdName] = fn
}

// **********************************************************************
// 実行
func (vm *VM) Run() error {
	// resultを受け取る準備
	if len(vm.AST.List) <= 1 {
		vm.Result = NewNil() // ASTが空ならnil, ASTが1つでも上書きを期待してnil
	} else {
		vm.Result = NewList(make([]Value, len(vm.AST.List))) // ASTが複数ならリスト
	}

	// ASTを順番に評価していく
	for i, v := range vm.AST.List {
		// 深さのリセット
		vm.Vars["vm_depth"] = NewNumber(0)
		vm.Vars["vm_depth_max"] = NewNumber(0)

		// 評価
		if ret, err := vm.Eval(v); err != nil {
			return err
		} else {
			// 結果を保存
			if vm.Result.IsList() {
				vm.Result.List[i] = ret
			} else {
				vm.Result = ret
			}
		}
	}

	return nil
}
