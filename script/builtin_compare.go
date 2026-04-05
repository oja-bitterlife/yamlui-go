package script

// **********************************************************************
// 比較系のbuiltin
var compareCmds = map[string]func(*VM, []Value) (Value, error){
	">":   grator,
	"<":   less,
	"==":  eq,
	"!=":  neq,
	">=":  ge,
	"<=":  le,
	"not": not,
	"and": and,
	"or":  or,
}

// ==================================================
// 比較演算
func grator(vm *VM, args []Value) (Value, error) {
	return binOp(vm, ">", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeString && arg1.Type == TypeString {
			return NewBool(arg0.Str > arg1.Str), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewBool(arg0.Num > arg1.Num), nil
		}
		return Value{}, binOpTypeError(">", arg0, arg1)
	})
}

func less(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "<", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeString && arg1.Type == TypeString {
			return NewBool(arg0.Str < arg1.Str), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewBool(arg0.Num < arg1.Num), nil
		}
		return Value{}, binOpTypeError("<", arg0, arg1)
	})
}

func eq(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "==", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeString && arg1.Type == TypeString {
			return NewBool(arg0.Str == arg1.Str), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewBool(arg0.Num == arg1.Num), nil
		}
		if arg0.Type == TypeBool && arg1.Type == TypeBool {
			return NewBool(arg0.Bool == arg1.Bool), nil
		}
		return Value{}, binOpTypeError("==", arg0, arg1)
	})
}

func neq(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "!=", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeString && arg1.Type == TypeString {
			return NewBool(arg0.Str != arg1.Str), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewBool(arg0.Num != arg1.Num), nil
		}
		if arg0.Type == TypeBool && arg1.Type == TypeBool {
			return NewBool(arg0.Bool != arg1.Bool), nil
		}
		return Value{}, binOpTypeError("!=: ", arg0, arg1)
	})
}

func ge(vm *VM, args []Value) (Value, error) {
	return binOp(vm, ">=", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeString && arg1.Type == TypeString {
			return NewBool(arg0.Str >= arg1.Str), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewBool(arg0.Num >= arg1.Num), nil
		}
		return Value{}, binOpTypeError(">=: ", arg0, arg1)
	})
}

func le(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "<=", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeString && arg1.Type == TypeString {
			return NewBool(arg0.Str <= arg1.Str), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewBool(arg0.Num <= arg1.Num), nil
		}
		return Value{}, binOpTypeError("<=: ", arg0, arg1)
	})
}

// ==================================================
// 論理演算
func not(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "not", args, func(vm *VM, arg0 Value) (Value, error) {
		// 0と0以外の数値を反転
		if arg0.Type == TypeNumber {
			if arg0.Num != 0 {
				return NewNumber(0), nil
			}
			return NewNumber(1), nil
		}
		// boolを反転
		if arg0.Type == TypeBool {
			return NewBool(!arg0.Bool), nil
		}

		// その他の値はエラー
		return Value{}, oneOpTypeError("not", arg0)
	})
}

func and(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "and", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.IsList() || arg1.IsList() {
			return Value{}, binOpTypeError("and", arg0, arg1)
		}

		return NewBool(arg0.ConvertBool().Bool && arg1.ConvertBool().Bool), nil
	})
}

func or(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "or", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.IsList() || arg1.IsList() {
			return Value{}, binOpTypeError("and", arg0, arg1)
		}

		return NewBool(arg0.ConvertBool().Bool || arg1.ConvertBool().Bool), nil
	})
}
