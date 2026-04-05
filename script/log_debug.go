//go:build debug

package script

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// ANSI カラーコード
const (
	bgInfo     = "\033[96;1m"
	bgError    = "\033[91;1m"
	bgFunc     = "\033[32m"
	colorGray  = "\033[97;2m"
	colorReset = "\033[0m"
)

type PrettyHandler struct {
	out io.Writer
}

func (h *PrettyHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	// レベルごとの色決定
	var levelStr, levelColor string
	switch r.Level {
	case slog.LevelError:
		levelStr, levelColor = "ERRO", bgError
	default:
		levelStr, levelColor = "INFO", bgInfo
	}

	// ソース情報の取得
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	source := fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)

	// フォーマット: [日時] [レベル] [ソース] | [メッセージ]
	// ※画像に合わせてセパレータやスペースを調整
	fmt.Fprintf(h.out, "%s%s%s %s%s%s %s%-20s %s| %s\n",
		colorGray, r.Time.Format("2006-01-02 15:04:05"), colorReset,
		levelColor, levelStr, colorReset,
		bgFunc, source, colorReset, r.Message,
	)
	return nil
}

// IFを満たすための実装
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *PrettyHandler) WithGroup(name string) slog.Handler       { return h }

var logger = slog.New(&PrettyHandler{out: os.Stderr})

func Log(msg string, args ...any) {
	pc, _, _, _ := runtime.Caller(1)
	r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(msg, args...), pc)
	_ = logger.Handler().Handle(context.Background(), r)
}

func LogErr(msg string, args ...any) error {
	fmtMsg := fmt.Sprintf(msg, args...)
	pc, _, _, _ := runtime.Caller(1)
	r := slog.NewRecord(time.Now(), slog.LevelError, fmtMsg, pc)
	_ = logger.Handler().Handle(context.Background(), r)
	return fmt.Errorf(fmtMsg)
}

func LogFatal(msg string, args ...any) {
	fmtMsg := fmt.Sprintf(msg, args...)
	pc, _, _, _ := runtime.Caller(1)
	r := slog.NewRecord(time.Now(), slog.LevelError, fmtMsg, pc)
	_ = logger.Handler().Handle(context.Background(), r)
	panic(fmtMsg)
}
