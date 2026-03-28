package script

import "errors"

func Itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append(buf, byte('0'+(n%10)))
		n /= 10
	}
	// 反転処理
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

func ValidationArgNum(cmdName string, argNum int, args []Value) error {
	if len(args) != argNum {
		argNumStr := Itoa(argNum)
		return errors.New(cmdName + ": requires exactly " + argNumStr + " arguments")
	}
	return nil
}
