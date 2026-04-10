//go:build debug

package script

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type PrettyHandler struct {
	out io.Writer
}

func SetLogWriter(w io.Writer) {
	logger = slog.New(&PrettyHandler{out: w})
}

// **********************************************************************
// ログのフォーマット
// ANSI カラーコード
const (
	bgInfo     = "\033[96;1m"
	bgWarn     = "\033[93;1m"
	bgError    = "\033[91;1m"
	bgFunc     = "\033[32m"
	colorGray  = "\033[97;2m"
	colorReset = "\033[0m"
)

func (h *PrettyHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	// レベルごとの色決定
	var levelStr, levelColor string
	switch r.Level {
	case slog.LevelError:
		levelStr, levelColor = "ERRO", bgError
	case slog.LevelWarn:
		levelStr, levelColor = "WARN", bgWarn
	default:
		levelStr, levelColor = "INFO", bgInfo
	}

	// ソース情報の取得
	source := ""
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		source = fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
	}

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

// **********************************************************************
// ログ出力
type LogLevel int

const (
	LogLevelInfo  LogLevel = LogLevel(slog.LevelInfo)
	LogLevelWarn  LogLevel = LogLevel(slog.LevelWarn)
	LogLevelError LogLevel = LogLevel(slog.LevelError)
)

// ==================================================
// 共通出力関数
func PrintLogCommon(level LogLevel, caller int, msg string, args ...any) string {
	fmtMsg := fmt.Sprintf(msg, args...)

	// 1. 呼び出し元のPCを取得
	// [呼び出し元] -> [Log/LogWarn] -> [PrintLogCommon]
	// なので、3階層上を指定する必要があります
	var pc uintptr
	var pcs [1]uintptr
	// runtime.Callers を使うのが一般的です (skip=3: PrintLogCommon, Log, 呼び出し元)
	runtime.Callers(caller, pcs[:])
	pc = pcs[0]

	r := slog.NewRecord(time.Now(), slog.Level(level), fmtMsg, pc)
	_ = logger.Handler().Handle(context.Background(), r)
	return fmtMsg
}

// ==================================================
// 各ログレベルの出力
func Log(msg string, args ...any) {
	PrintLogCommon(LogLevelInfo, 3, msg, args...)
}

func LogWarn(msg string, args ...any) {
	PrintLogCommon(LogLevelWarn, 3, msg, args...)
}

func LogErr(msg string, args ...any) error {
	fmtMsg := PrintLogCommon(LogLevelError, 3, msg, args...)
	return errors.New(fmtMsg)
}

func LogFatal(msg string, args ...any) {
	fmtMsg := PrintLogCommon(LogLevelError, 3, msg, args...)
	panic(fmtMsg)
}
