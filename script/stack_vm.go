package script

import "fmt"

// HostFunction は外部（UI側）で定義される関数の型
type HostFunction func(args []Value) Value

// FunctionEntry は引数の数と実行体をセットにしたもの
type FunctionEntry struct {
	Arity int
	Exec  HostFunction
}

type StackVM struct {
	stack []Value
	vars  map[string]float64 // @property 用のストレージ
}

func NewStackVM() *StackVM {
	return &StackVM{
		stack: make([]Value, 0),
		vars:  make(map[string]float64),
	}
}

// Execute は命令列と、その時々で必要な関数テーブルを受け取って実行します
func (vm *StackVM) Execute(insts []Instruction, funcs map[string]FunctionEntry) {
	for _, ins := range insts {
		switch ins.Op {
		case OpPushNum, OpPushBool, OpPushStr, OpPushSym, OpPushProp:
			vm.push(ins.Val)

		case OpCall:
			vm.dispatch(ins.Val.Str, funcs)
		}
	}
}

// dispatch は関数名を見て、外部関数か組み込み関数かを振り分けます
func (vm *StackVM) dispatch(name string, funcs map[string]FunctionEntry) {
	// 1. 外部注入された関数にあるか確認
	if fn, ok := funcs[name]; ok {
		args := make([]Value, fn.Arity)
		// スタックなので「逆順」に取り出して引数スライスを作る
		for i := fn.Arity - 1; i >= 0; i-- {
			args[i] = vm.pop()
		}
		vm.push(fn.Exec(args))
		return
	}

	// 2. 無ければ VM 内部の組み込みロジック（_.set など）へ
	vm.executeBuiltin(name)
}

func (vm *StackVM) executeBuiltin(name string) {
	switch name {
	case "_.set":
		val := vm.pop().Number
		prop := vm.pop().Str // OpPushProp で積まれた名前
		vm.vars[prop] = val
		// デバッグ用
		fmt.Printf("[StackVM] @%s = %g\n", prop, val)

	default:
		fmt.Printf("[StackVM Error] Unknown function: %s\n", name)
	}
}

func (vm *StackVM) push(v Value) {
	vm.stack = append(vm.stack, v)
}

func (vm *StackVM) pop() Value {
	if len(vm.stack) == 0 {
		return Value{Type: TypeNil}
	}
	v := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return v
}
