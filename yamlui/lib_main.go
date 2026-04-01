package yamlui

import (
	"errors"
	"slices"

	"github.com/oja-bitterlife/yamlui-go/script"
)

type YAMLUI struct {
	// UIツリー構築用
	Root       *UIBase
	remapFuncs map[string]func(type_ string, parent *UIBase, data map[string]script.Value) (*UIBase, error)

	// Updateの時に使うもの
	frame       int // システム時間
	EventQueue  []string
	updateQueue []UpdateQueueItem

	// Drawの時に使うもの
	Screen    Area
	drawQueue []DrawQueueItem
}

func NewYAMLUI() *YAMLUI {
	return &YAMLUI{
		Root:       NewUIBase("Root"),
		remapFuncs: make(map[string]func(type_ string, parent *UIBase, data map[string]script.Value) (*UIBase, error)),
	}
}

// **********************************************************************
// mapの解析
// ==================================================
// mapを解析して必要なデータを構造体に流し込むためのインターフェース
func (self *YAMLUI) RegisterRemap(typeName string, fn func(type_ string, parent *UIBase, data map[string]script.Value) (*UIBase, error)) {
	self.remapFuncs[typeName] = fn
}

// UIを構築するためのインターフェース
type UIComponent interface {
	GetBase() *UIBase
}

type UIComponentFactory[T UIComponent] func(componentName string, parent *UIBase, data map[string]script.Value) T

func UIBuilder[T UIComponent](
	lib *YAMLUI,
	componentName string,
	factory UIComponentFactory[T],
	onCreated func(T)) {

	lib.RegisterRemap(componentName,
		// クロージャでUIComponentを構築
		func(type_ string, parent *UIBase, data map[string]script.Value) (*UIBase, error) {
			// Factoryで構築
			component := factory(componentName, parent, data)

			// onCreatedがnilでなければ呼び出す
			if onCreated != nil {
				onCreated(component)
			}

			// インターフェース経由で Base を取り出して返す
			return component.GetBase(), nil
		})
}

// ==================================================
// map解析
// Mapの値を構造体に流し込むためのヘルパー関数
func PropStr(data map[string]script.Value, key string, def string) string {
	value, ok := data[key]
	if !ok || value.Type != script.TypeString {
		return def
	}
	return value.Str
}

func PropNum(data map[string]script.Value, key string, def float64) float64 {
	value, ok := data[key]
	if !ok || value.Type != script.TypeNumber {
		return def
	}
	return value.Num
}

func PropINum(data map[string]script.Value, key string, def int) int {
	value, ok := data[key]
	if !ok || value.Type != script.TypeNumber {
		return def
	}
	return int(value.Num)
}

func PropBool(data map[string]script.Value, key string, def bool) bool {
	value, ok := data[key]
	if !ok || value.Type != script.TypeBool {
		return def
	}
	return value.Bool
}

// **********************************************************************
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

// **********************************************************************
// 呼び出し口
// Update
func (self *YAMLUI) Update(frame int) []error {
	errorList := []error{}

	self.frame = frame

	// updateQueueをクリア
	self.updateQueue = []UpdateQueueItem{}

	// 更新コンテキストを作成してUpdateTreeを呼び出す
	// ----------------------------------------
	ctx := NewUpdateContext(self, self.Root, nil)
	if err := self.Root.RecUpdateTree(0, ctx); err != nil {
		errorList = append(errorList, err...)
	}

	// drawQueueに溜まった描画命令を実行する
	// ----------------------------------------
	// Z順でソートする
	slices.SortStableFunc(self.updateQueue, func(a, b UpdateQueueItem) int {
		return a.z - b.z
	})

	// UpdateCountが0のときはInitとみなしてOnInitを呼び出す
	for _, item := range self.updateQueue {
		uiBase := item.ctx.Base
		if uiBase.UpdateCount == 0 && uiBase.onInitIF != nil {
			if err := uiBase.onInitIF.OnInit(item.ctx); err != nil {
				errorList = append(errorList, err)
			}
		}
	}

	// イベントの処理
	self.ProcessEvents()

	// ソートされたqueueを順番に実行する
	for _, item := range self.updateQueue {
		item.UpdateIF.Update(item.ctx)

		// UIの更新後スクリプトがあれば走らせる
		uiBase := item.ctx.Base
		if uiBase.script != nil {
			uiBase.storeToVM(uiBase.script.GetVM())
			if _, err := uiBase.script.Run(); err != nil {
				errorList = append(errorList, err)
			}
			uiBase.loadFromVM(uiBase.script.GetVM())
		}

		// 実行カウントを進める
		uiBase.UpdateCount++
	}

	return errorList
}

// ==================================================
// Draw
func (self *YAMLUI) Draw(screen Area) {
	self.Screen = screen

	// drawQueueをクリア
	self.drawQueue = []DrawQueueItem{}

	// 描画コンテキストを作成してDrawTreeを呼び出す
	// ----------------------------------------
	ctx := NewDrawContext(self, self.Root, nil, screen)
	self.Root.RecDrawTree(0, 0, 0, ctx)

	// drawQueueに溜まった描画命令を実行する
	// ----------------------------------------
	// Z順でソートする
	slices.SortStableFunc(self.drawQueue, func(a, b DrawQueueItem) int {
		return a.z - b.z
	})
	// ソートされたqueueを順番に実行する
	for _, item := range self.drawQueue {
		item.drawIF.Draw(item.x, item.y, item.ctx)
	}
}
