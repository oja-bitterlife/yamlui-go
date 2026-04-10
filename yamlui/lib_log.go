package yamlui

import (
	"errors"

	"github.com/oja-bitterlife/yamlui-go/script"
)

func (lib *YAMLUI) Log(msg string, args ...any) {
	script.PrintLogCommon(script.LogLevelInfo, 4, msg, args...)
}
func (lib *YAMLUI) LogWarn(msg string, args ...any) {
	script.PrintLogCommon(script.LogLevelWarn, 4, msg, args...)
}
func (lib *YAMLUI) LogErr(msg string, args ...any) error {
	fmtMsg := script.PrintLogCommon(script.LogLevelError, 4, msg, args...)
	return errors.New(fmtMsg)
}
