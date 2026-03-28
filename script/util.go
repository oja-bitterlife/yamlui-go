package script

import (
	"errors"
	"strconv"
)

func ValidationArgNum(cmdName string, argNum int, args []Value) error {
	if len(args) != argNum {
		argNumStr := strconv.Itoa(argNum)
		return errors.New(cmdName + ": requires exactly " + argNumStr + " arguments")
	}
	return nil
}
