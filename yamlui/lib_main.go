package yamlui

import (
	"errors"
	"path"
	"slices"

	"github.com/oja-bitterlife/yamlui-go/script"
)

type YAMLUI struct {
	// UIツリー構築用
	root   *UIBase
	refObj map[string]UIComponent[*UIBase] // Component生成用リファレンスオブジェクト

	// Updateの時に使うもの
	SystemFrame int // システム時間
	eventQueue  []string
	updateQueue []UpdateQueueItem

	// Drawの時に使うもの
	Screen    Area
	drawQueue []DrawQueueItem
}

func NewYAMLUI() *YAMLUI {
	return &YAMLUI{
		root:   NewUIBase(),
		refObj: make(map[string]UIComponent[*UIBase]),
	}
}

func (self *YAMLUI) FindByID(id string) *UIBase {
	return self.root.FindChildByID(id)
}

// **********************************************************************
// UIのJSONをUnmarshalしたValueMap(map[string]Value)のTypeごとにインスタンスを割り当て
// ==================================================
type UIComponent[T any] interface {
	GetUIBase() *UIBase
	Clone() UIComponent[*UIBase]
	Setup(type_ string, data script.ValueMap) error
}

// LoadのときにTypeを見て登録されたUICloneableからUIBaseを複製して構築するインターフェース
func (self *YAMLUI) UIBuild(type_ string, refObj UIComponent[*UIBase]) {
	self.refObj[type_] = refObj
}

// ==================================================
// map解析
// Mapの値を構造体に流し込むためのヘルパー関数
// **********************************************************************
// UITreeの構築（再帰的に子要素も構築）
func (self *YAMLUI) Load(data []byte) error {
	// JSONからValueを経由してUIBaseを再構築する
	var value script.Value
	if err := value.UnmarshalJSON(data); err != nil {
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
		return errors.New("Expected top-level Value to be List or Map, got " + value.Type.String())
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
			return errors.New("Invalid pattern in UIBuild: " + pattern + ": " + err.Error())
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
		ui.LoadFromValue(value) // プロパティを流し込む

		// Setup関数でさらに細かい構築を行う
		if err := component.Setup(type_, value.Map); err != nil {
			return err
		}
	} else {
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

// **********************************************************************
// 呼び出し口
// Update
func (self *YAMLUI) Update(frame int) []error {
	errorList := []error{}

	self.SystemFrame = frame

	// updateQueueをクリア
	self.updateQueue = []UpdateQueueItem{}

	// 更新コンテキストを作成してUpdateTreeを呼び出す
	// ----------------------------------------
	if err := self.root.RecUpdateTree(self, 0); err != nil {
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
		uiBase := item.UpdateIF.GetUIBase()
		if uiBase.UpdateCount == 0 && uiBase.onInitIF != nil {
			if err := uiBase.onInitIF.OnInit(self); err != nil {
				errorList = append(errorList, err)
			}
		}
	}

	// イベントをUpdateの前に処理し、Updateにイベントを通知する
	self.ProcessEvents()
	self.ClearEvents() // 終わったらイベントキューをクリア

	// ソートされたqueueを順番に実行する
	for _, item := range self.updateQueue {
		uiBase := item.UpdateIF.GetUIBase()
		uiBase.Action = "" // Actionをクリアしておく

		// Updateを呼び出す
		// ----------------------------------------
		item.UpdateIF.Update(self, item.Events)

		// Update後スクリプトがあれば走らせる
		// ----------------------------------------
		if uiBase.script != nil {
			// スクリプトを実行する前に、UIBaseのプロパティをVMに保存しておく
			uiBase.storeToVM(uiBase.script.GetVM())

			// スクリプトを実行
			if err := uiBase.script.Run(); err != nil {
				errorList = append(errorList, err)
			}

			// スクリプトを実行した後に、VMからUIBaseのプロパティを更新する
			uiBase.loadFromVM(uiBase.script.GetVM())
		}

		// Update後処理
		// ----------------------------------------
		// Updateの実行後にイベントが発生(Action!="")していればキューに追加
		if uiBase.Action != "" {
			self.AddEvent(uiBase.Action)
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
	ctx := NewDrawContext(self, self.root, nil, screen)
	self.root.RecDrawTree(0, 0, 0, ctx)

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
