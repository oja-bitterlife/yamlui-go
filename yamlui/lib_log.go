package yamlui

import "github.com/oja-bitterlife/yamlui-go/script"

func (lib *YAMLUI) Log(msg string, args ...any) {
	script.Log(msg, args...)
}
func (lib *YAMLUI) LogWarn(msg string, args ...any) {
	script.LogWarn(msg, args...)
}
func (lib *YAMLUI) LogErr(msg string, args ...any) error {
	return script.LogErr(msg, args...)
}
