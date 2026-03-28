package yamlui

import (
	dbm "github.com/oja-bitterlife/double-mappering-go"
	"github.com/oja-bitterlife/yamlui-go/builtin"
	"github.com/oja-bitterlife/yamlui-go/script"
)

type YamlUI[T any] struct {
	dbm *dbm.DoubleMappering[T] // データベースマッピング
	vm  *script.VM              // 自作されたLisp処理系
}

func NewYamlUI[T any](
	// データのシリアライズとデシリアライズの関数
	// updateはトランザクション内で行われる
	marshal func(*T) ([]byte, error),
	unmarshal func([]byte) (*T, error),
) *YamlUI[T] {
	dbm := dbm.New(marshal, unmarshal)

	return &YamlUI[T]{
		dbm: dbm,
		vm:  builtin.NewWithBuiltin(),
	}
}

func (ui *YamlUI[T]) Init() {
	// ui.init(ui)
}

func (ui *YamlUI[T]) Update() {
	// ui.update(ui)
}

func (u *YamlUI[T]) Draw() {
	// ui.draw(ui)
}
