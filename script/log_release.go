//go:build !debug

package script

import "errors"

// 完全に消えることを期待
func Log(msg string, args ...any) {}
func LogErr(msg string, args ...any) error {
	return errors.New(msg)
}
func LogFatal(msg string, args ...any) {
	panic(msg)
}
