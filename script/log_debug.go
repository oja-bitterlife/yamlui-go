//go:build debug

package script

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	AddSource: true,
}))

func Log(msg string, args ...any) {
	logger.Info(msg, args...)
}

func LogErr(msg string, args ...any) error {
	fmtMsg := fmt.Sprintf(msg, args...)

	// 呼び出し元（LogFatal を呼んだ場所）
	pc, _, _, _ := runtime.Caller(1)
	r := slog.NewRecord(time.Now(), slog.LevelError, fmtMsg, pc)

	// ハンドラを直接叩くことで、slog 内部の runtime.Caller(深さ固定) をバイパスする
	_ = logger.Handler().Handle(context.Background(), r)

	return fmt.Errorf(fmtMsg)
}

func LogFatal(msg string, args ...any) {
	fmtMsg := fmt.Sprintf(msg, args...)

	// 呼び出し元（LogFatal を呼んだ場所）
	pc, _, _, _ := runtime.Caller(1)
	r := slog.NewRecord(time.Now(), slog.LevelError, fmtMsg, pc)

	// ハンドラを直接叩くことで、slog 内部の runtime.Caller(深さ固定) をバイパスする
	_ = logger.Handler().Handle(context.Background(), r)

	panic(fmtMsg)
}
