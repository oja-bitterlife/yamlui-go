package yamlui

import (
	"errors"
	"log/slog"

	"github.com/oja-bitterlife/yamlui-go/script"
)

func (lib *YAMLUI) Log(msg string, args ...any) {
	script.PrintLogCommon(slog.LevelInfo, 4, msg, args...)
}
func (lib *YAMLUI) LogWarn(msg string, args ...any) {
	script.PrintLogCommon(slog.LevelWarn, 4, msg, args...)
}
func (lib *YAMLUI) LogErr(msg string, args ...any) error {
	fmtMsg := script.PrintLogCommon(slog.LevelError, 4, msg, args...)
	return errors.New(fmtMsg)
}
