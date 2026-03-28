package yamlui

import (
	"encoding/json"

	dbm "github.com/oja-bitterlife/double-mappering-go"
	"github.com/oja-bitterlife/yamlui-go/builtin"
	"github.com/oja-bitterlife/yamlui-go/script"
)

type YamlUI[T any] struct {
	dbm *dbm.DoubleMappering[T] // データベースマッピング
	vm  *script.VM              // 自作されたLisp処理系
}

func NewYamlUI[T any]() *YamlUI[T] {
	dbm := dbm.New(
		func(data *T) ([]byte, error) {
			return json.Marshal(data)
		},
		func(b []byte) (*T, error) {
			var data T
			err := json.Unmarshal(b, &data)
			return &data, err
		},
	)

	return &YamlUI[T]{
		dbm: dbm,
		vm:  builtin.NewWithBuiltin(),
	}
}
