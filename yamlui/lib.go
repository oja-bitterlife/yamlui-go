package yamlui

import (
	"errors"
)

type YAMLUI struct {
	Root       *UIBase
	remapFuncs map[string]func(*UIBase, map[string]any) (*UIBase, error)
}

func NewYAMLUI() *YAMLUI {
	return &YAMLUI{
		Root:       NewUIBase("Root"),
		remapFuncs: make(map[string]func(*UIBase, map[string]any) (*UIBase, error)),
	}
}

// **********************************************************************
// mapの解析
// ==================================================
// mapを解析して必要なデータを構造体に流し込むためのインターフェース
func (self *YAMLUI) RegisterRemap(typeName string, fn func(*UIBase, map[string]any) (*UIBase, error)) {
	self.remapFuncs[typeName] = fn
}

// ==================================================
// map解析
// Mapの値を構造体に流し込むためのヘルパー関数
func propStr(data map[string]any, key string, def string) string {
	if v, ok := data[key].(string); ok {
		return v
	}
	return def
}

func propNum(data map[string]any, key string, def float64) float64 {
	if v, ok := data[key].(float64); ok {
		return v
	}
	return def
}
func propBool(data map[string]any, key string, def bool) bool {
	if v, ok := data[key].(bool); ok {
		return v
	}
	return def
}

// ==================================================
// UIコンポーネントの構築
func (self *YAMLUI) BuildUI(data map[string]any) (*UIBase, error) {
	Type := propStr(data, "Type", "area") // Type を取得（デフォルトは "area"）
	ui := NewUIBase(Type)                 // データ格納用

	if ID, ok := data["ID"].(string); ok {
		ui.ID = ID
	}

	// 共通プロパティの取り出し
	ui.IsEnable = propBool(data, "IsEnable", ui.IsEnable)
	ui.IsAbs = propBool(data, "IsAbs", ui.IsAbs)
	ui.X = propNum(data, "X", ui.X)
	ui.Y = propNum(data, "Y", ui.Y)
	ui.W = propNum(data, "W", ui.W)
	ui.H = propNum(data, "H", ui.H)

	ui.IsVisible = propBool(data, "IsVisible", ui.IsVisible)
	ui.Text = propStr(data, "Text", ui.Text)
	ui.Color = propStr(data, "Color", ui.Color)

	ui.SelectNo = propNum(data, "SelectNo", ui.SelectNo)
	ui.SelGridX = propNum(data, "SelGridX", ui.SelGridX)
	ui.Action = propStr(data, "Action", ui.Action)

	// ユーザー定義コンポーネントを生成
	if remapFunc, ok := self.remapFuncs[Type]; ok {
		// 登録されたリマップ関数でUIを構築
		if remapUI, err := remapFunc(ui, data); err == nil {
			remapUI.CopyProp(ui) // 共通プロパティをコピー
			return remapUI, nil
		} else {
			return nil, err
		}
	} else {
		// 登録されたリマップ関数がない場合は、基本的なUIを返す
		return ui, nil
	}
}

// ==================================================
// UITreeの構築（再帰的に子要素も構築）
func (self *YAMLUI) Load(data any) error {
	switch v := data.(type) {
	case []any:
		return errors.New("invalid data format: expected a single map, but got an array")
	case []map[string]any:
		// ルートが配列 [] の場合
		for _, item := range v {
			if err := self.load(self.Root, item); err != nil {
				return err
			}
		}
	case map[string]any:
		// ルートが単一オブジェクト {} の場合
		return self.load(self.Root, v)
	default:
		return errors.New("invalid data format: expected array or map")
	}
	return nil
}

func (self *YAMLUI) load(parent *UIBase, data map[string]any) error {
	// 現在のノードを構築
	ui, err := self.BuildUI(data)
	if err != nil {
		return err
	}
	parent.AddChild(ui) // 親ノードに追加

	// ここで再帰的に UI を構築するロジックを入れる
	children, ok := data["children"].([]any)
	if !ok {
		return nil // children がない場合は終了
	}
	for _, c := range children {
		if m, ok := c.(map[string]any); ok {
			if err := self.load(ui, m); err != nil {
				return err
			}
		}
	}
	return nil
}
