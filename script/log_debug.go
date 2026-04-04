//go:build debug

package script

import (
	"fmt"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	AddSource: true,
}))

func Log(msg string, args ...any) {
	logger.Info(msg, args...)
}

func LogErr(msg string, args ...any) error {
	logger.Error(msg, args...)
	return fmt.Errorf(msg, args...)
}

func LogFatal(msg string, args ...any) {
	logger.Error(msg, args...)
	panic(fmt.Sprintf(msg, args...))
}
