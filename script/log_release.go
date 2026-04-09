//go:build !debug

package script

import (
	"errors"
	"io"
	"log/slog"
)

// 完全に消えることを期待
// ==================================================
func SetLogWriter(w io.Writer)                                                       {}
func PrintLogCommon(level slog.Level, callDepth int, msg string, args ...any) string { return "" }

// ----------------------------------------
func Log(msg string, args ...any)     {}
func LogWarn(msg string, args ...any) {}
func LogErr(msg string, args ...any) error {
	return errors.New(msg)
}
func LogFatal(msg string, args ...any) {
	panic(msg)
}
