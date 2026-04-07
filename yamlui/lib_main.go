package yamlui

import (
	"path"

	"github.com/oja-bitterlife/yamlui-go/script"
)

const (
	EVENT_USING_MAX = 8  // 同時に発生するイベントの最大数
	UI_USING_MAX    = 32 // 同時に使用するUIの最大数
)

type YAMLUI struct {
	// UIツリー構築用
	root   *UIBase
	refObj map[string]UIComponent[*UIBase] // Component生成用リファレンスオブジェクト

	// Updateの時に使うもの
	SystemFrame int // システム時間
	eventQueue  []EventQueueItem
	updateQueue []UpdateQueueItem

	// Drawの時に使うもの
	Screen    Area
	drawQueue []DrawQueueItem
}

func NewYAMLUI() *YAMLUI {
	return &YAMLUI{
		root:   NewUIBase(),
		refObj: make(map[string]UIComponent[*UIBase]),

		eventQueue:  make([]EventQueueItem, 0, EVENT_USING_MAX),
		updateQueue: make([]UpdateQueueItem, 0, UI_USING_MAX),

		drawQueue: make([]DrawQueueItem, 0, UI_USING_MAX),
	}
}

// ==================================================
// Getter/Setter
func (self *YAMLUI) FindByID(id string) *UIBase {
	return self.root.FindChildByID(id)
}

func (self *YAMLUI) GetRoot() *UIBase {
	return self.root
}

// デバッグ用。ないとログもままならないので
func (self *YAMLUI) GetEvents() []EventQueueItem {
	return self.eventQueue
}

// **********************************************************************
// UIのJSONをUnmarshalしたValueMap(map[string]Value)のTypeごとにインスタンスを割り当て
// ==================================================
type UIComponent[T any] interface {
	GetUIBase() *UIBase
	Clone() UICloned
	Setup(type_ string, data script.ValueMap) error
}

// 簡易アクセス用
type UICloned UIComponent[*UIBase]

// LoadのときにTypeを見て登録されたUICloneableからUIBaseを複製して構築するインターフェース
func (self *YAMLUI) UIBuild(type_ string, refObj UIComponent[*UIBase]) {
	self.refObj[type_] = refObj
}

// **********************************************************************
// UITreeの構築（再帰的に子要素も構築）
func (self *YAMLUI) Load(valueJson []byte) error {
	// JSONからValueを経由してUIBaseを再構築する
	value, err := script.NewFromValueJSON(valueJson)
	if err != nil {
		return err
	}

	// Loadで再帰的にUIを構築する
	self.root.children = []*UIBase{} // 既存の子要素をクリア

	switch value.Type {
	case script.TypeLitList:
		// 最初が配列なら最初ですべてを子要素として追加する
		for _, item := range value.List {
			err := self.load(self.root, item)
			if err != nil {
				return err
			}
		}
		return nil
	case script.TypeLitMap:
		// 最初がMapなら普通に登録していく
		err := self.load(self.root, value)
		if err != nil {
			return err
		}
	default:
		return script.LogErr("Expected top-level Value to be List or Map, got " + value.Type.String())
	}

	return nil
}

// 現在のノードを構築
func (self *YAMLUI) load(parent *UIBase, value script.Value) error {
	// Typeを見て、リマップ関数があればそれで構築
	type_ := value.Map["Type"].Str

	// path.Matchを使ってtype_を曖昧にマッチさせる
	matchObj := UIComponent[*UIBase](nil)
	for pattern, refObj := range self.refObj {
		matched, err := path.Match(pattern, type_)
		if err != nil {
			return script.LogErr("Invalid pattern in UIBuild: " + pattern + ": " + err.Error())
		}
		if matched {
			matchObj = refObj
			break
		}
	}

	var ui *UIBase
	if matchObj != nil {
		// 登録されたUICloneableからUIを複製して構築
		component := matchObj.Clone()
		ui = component.GetUIBase()

		// ValueからUIBaseに流し込む
		if err := ui.LoadFromValue(value); err != nil {
			return err
		}

		// Setup関数でさらに細かい構築を行う
		if err := component.Setup(type_, value.Map); err != nil {
			return err
		}
	} else {
		// エラー終了せず警告で進めてみる
		script.LogWarn("Unregister Type: " + type_)

		// 登録されたUICloneableがない場合は、基本的なUIを構築
		ui = NewUIBase()
		ui.LoadFromValue(value) // プロパティを流し込む
	}

	// 構築したUIを親に追加
	parent.AddChild(ui)

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
