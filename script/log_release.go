//go:build !debug

package script

import (
	"errors"
	"io"
)

// 完全に消えることを期待
// ==================================================
type LogLevel int

const (
	// slogに合わせておく
	LogLevelInfo  LogLevel = 0
	LogLevelWarn  LogLevel = 4
	LogLevelError LogLevel = 8
)

func SetLogWriter(w io.Writer)                                                     {}
func PrintLogCommon(level LogLevel, callDepth int, msg string, args ...any) string { return "" }

// ----------------------------------------
func Log(msg string, args ...any)     {}
func LogWarn(msg string, args ...any) {}
func LogErr(msg string, args ...any) error {
	return errors.New(msg)
}
func LogFatal(msg string, args ...any) {
	panic(msg)
}
