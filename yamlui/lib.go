package yamlui

import (
	"errors"

	"github.com/oja-bitterlife/yamlui-go/script"
)

type YAMLUI struct {
	Root       *UIBase
	remapFuncs map[string]func(string, *UIBase, map[string]script.Value) (*UIBase, error)
}

func NewYAMLUI() *YAMLUI {
	return &YAMLUI{
		Root:       NewUIBase("Root"),
		remapFuncs: make(map[string]func(string, *UIBase, map[string]script.Value) (*UIBase, error)),
	}
}

// **********************************************************************
// mapの解析
// ==================================================
// mapを解析して必要なデータを構造体に流し込むためのインターフェース
func (self *YAMLUI) RegisterRemap(typeName string, fn func(string, *UIBase, map[string]script.Value) (*UIBase, error)) {
	self.remapFuncs[typeName] = fn
}

// ==================================================
// map解析
// Mapの値を構造体に流し込むためのヘルパー関数
func propStr(data map[string]script.Value, key string, def string) string {
	value, ok := data[key]
	if !ok || value.Type != script.TypeString {
		return def
	}
	return value.Str
}

func propNum(data map[string]script.Value, key string, def float64) float64 {
	value, ok := data[key]
	if !ok || value.Type != script.TypeNumber {
		return def
	}
	return value.Num
}

func propINum(data map[string]script.Value, key string, def int) int {
	value, ok := data[key]
	if !ok || value.Type != script.TypeNumber {
		return def
	}
	return int(value.Num)
}

func propBool(data map[string]script.Value, key string, def bool) bool {
	value, ok := data[key]
	if !ok || value.Type == script.TypeNumber {
		return def
	}
	return def
}

// ==================================================
// UITreeの構築（再帰的に子要素も構築）
func (self *YAMLUI) Load(data []byte) error {
	// JSONからValueを経由してUIBaseを再構築する
	var value script.Value
	if err := value.UnmarshalJSON(data); err != nil {
		return err
	}

	// Loadで再帰的にUIを構築する
	self.Root.children = []*UIBase{} // 既存の子要素をクリア

	switch value.Type {
	case script.TypeLitList:
		// 最初が配列なら最初ですべてを子要素として追加する
		for _, item := range value.List {
			err := self.load(self.Root, item)
			if err != nil {
				return err
			}
		}
		return nil
	case script.TypeLitMap:
		// 最初がMapなら普通に登録していく
		err := self.load(self.Root, value)
		if err != nil {
			return err
		}
	default:
		return errors.New("Expected top-level Value to be List or Map, got " + value.Type.String())
	}

	return nil
}

func (self *YAMLUI) load(parent *UIBase, value script.Value) error {
	// 現在のノードを構築
	var ui *UIBase
	var err error

	// Typeを見て、リマップ関数があればそれで構築
	type_ := value.Map["Type"].Str
	if remapFunc, ok := self.remapFuncs[type_]; ok {
		// 登録されたリマップ関数でUIを構築
		ui, err = remapFunc(type_, parent, value.Map)
		if err != nil {
			return err
		}
	} else {
		// 登録されたリマップ関数がない場合は、基本的なUIを構築
		ui = NewUIBase(type_)
	}

	ui.LoadFromValue(value) // プロパティを流し込む
	parent.AddChild(ui)     // 親ノードに追加

	// ここで再帰的に UI を構築するロジックを入れる
	children, ok := value.Map["children"]
	if !ok {
		return nil // children がない場合は終了
	}
	for _, childValue := range children.List {
		// 再帰的に子ノードを構築
		err := self.load(ui, childValue)
		if err != nil {
			return err
		}
	}

	return nil
}
