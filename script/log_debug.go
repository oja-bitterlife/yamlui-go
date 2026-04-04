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
	bgInfo     = "\033[42;97m" // 緑背景・黒文字
	bgDebug    = "\033[46;97m" // シアン背景・黒文字
	bgWarn     = "\033[43;97m" // 黄背景・黒文字
	bgError    = "\033[41;97m" // 赤背景・白文字
	colorGray  = "\033[90m"    // 灰色（日時用）
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
	case slog.LevelDebug:
		levelStr, levelColor = "DEBUG", bgDebug
	case slog.LevelWarn:
		levelStr, levelColor = "WARN ", bgWarn
	case slog.LevelError:
		levelStr, levelColor = "ERROR", bgError
	default:
		levelStr, levelColor = "INFO ", bgInfo
	}

	// ソース情報の取得
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	source := fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)

	// フォーマット: [日時] [レベル] [ソース] | [メッセージ]
	// ※画像に合わせてセパレータやスペースを調整
	fmt.Fprintf(h.out, "%s%s%s %s%s%s %-20s | %s\n",
		colorGray, r.Time.Format("2006-01-02 15:04:05"), colorReset,
		levelColor, levelStr, colorReset,
		source, r.Message,
	)
	return nil
}

// 既存のメソッドは変更なし
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

func LogWarn(msg string, args ...any) {
	pc, _, _, _ := runtime.Caller(1)
	r := slog.NewRecord(time.Now(), slog.LevelWarn, fmt.Sprintf(msg, args...), pc)
	_ = logger.Handler().Handle(context.Background(), r)
}

func LogDebug(msg string, args ...any) {
	pc, _, _, _ := runtime.Caller(1)
	r := slog.NewRecord(time.Now(), slog.LevelDebug, fmt.Sprintf(msg, args...), pc)
	_ = logger.Handler().Handle(context.Background(), r)
}

func LogFatal(msg string, args ...any) {
	fmtMsg := fmt.Sprintf(msg, args...)
	pc, _, _, _ := runtime.Caller(1)
	r := slog.NewRecord(time.Now(), slog.LevelError, fmtMsg, pc)
	_ = logger.Handler().Handle(context.Background(), r)
	panic(fmtMsg)
}
