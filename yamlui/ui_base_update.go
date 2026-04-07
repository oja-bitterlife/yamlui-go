package yamlui

import (
	"slices"
)

// **********************************************************************
// Updateのインターフェース
// ==================================================
// 直接DrawIFを呼び出すのではなく、UpdateTreeの中でUpdateQueueItemにしてキューに入れる
type UpdateQueueItem struct {
	UpdateIF UpdateIF
	z        int
}

// ==================================================
// OnInit
type OnInitIF interface {
	OnInit(lib *YAMLUI) (string, error)
	GetUIBase() *UIBase
}

func (self *UIBase) SetOnInitIF(onInitIF OnInitIF) {
	self.onInitIF = onInitIF
}

// ==================================================
// Update
type UpdateIF interface {
	Update(lib *YAMLUI, events []string) (string, error)
	GetUIBase() *UIBase
}

type UpdateTreeIF interface {
	UpdateTree(lib *YAMLUI, z int) []error
}

func (self *UIBase) SetUpdateIF(updateIF UpdateIF) {
	self.updateIF = updateIF
}

func (self *UIBase) SetUpdateTreeIF(updateTreeIF UpdateTreeIF) {
	self.updateTreeIF = updateTreeIF
}

// **********************************************************************
// 呼び出し口
func (self *YAMLUI) Update(frame int) []error {
	errorList := []error{}

	self.SystemFrame = frame

	// updateQueueをクリア
	self.updateQueue = self.updateQueue[:0]

	// 更新コンテキストを作成してUpdateTreeを呼び出す
	// ----------------------------------------
	if err := self.root.recUpdateTree(self, 0); err != nil {
		errorList = append(errorList, err...)
	}

	// updateQueueに溜まったUpdateを実行する
	// ----------------------------------------
	// Z順でソートする
	slices.SortStableFunc(self.updateQueue, func(a, b UpdateQueueItem) int {
		return a.z - b.z
	})

	// UpdateCountが0のときはInitとみなしてOnInitを呼び出す
	for _, item := range self.updateQueue {
		uiBase := item.UpdateIF.GetUIBase()

		// OnInitを呼び出す
		// ----------------------------------------
		if uiBase.UpdateCount == 0 && uiBase.onInitIF != nil {
			event, err := uiBase.onInitIF.OnInit(self)
			if err != nil {
				errorList = append(errorList, err)
			}

			// OnIniteの実行後にイベントが発生(event!="")していればキューに追加
			if event != "" {
				self.AddEvent(event)
			}
		}
	}

	// イベントをUpdateの前に処理し、Updateにイベントを通知する
	self.ProcessEvents()
	self.eventQueue = self.eventQueue[:0] // イベントはここでクリア

	// ここ以降でself.eventQueueに入るイベントはUpateで入るイベント

	// ソートされたqueueを順番に実行する
	for qi, item := range self.updateQueue {
		uiBase := item.UpdateIF.GetUIBase()

		// スクリプトがあれば走らせる
		// ----------------------------------------
		if uiBase.HasScript() {
			// スクリプトを実行する前に、UIBaseのプロパティをVMに保存しておく
			uiBase.storeToVM()

			// スクリプトを実行
			if err := uiBase.script.Run(); err != nil {
				errorList = append(errorList, err)
			}

			// スクリプトを実行した後に、VMからUIBaseのプロパティを更新する
			uiBase.loadFromVM()
		}

		// Update
		// ----------------------------------------
		// Updateに渡すeventの回収
		events := []string{}
		for _, e := range self.eventQueue {
			if e.updateQueueNo == qi {
				events = append(events, e.event)
			}
		}

		// Updateを呼び出す
		event, err := item.UpdateIF.Update(self, events)
		if err != nil {
			errorList = append(errorList, err)
		}

		// Updateの実行後にイベントが発生(event!="")していればキューに追加
		if event != "" {
			self.AddEvent(event)
		}
	}

	// 実行カウントを進める
	// ----------------------------------------
	for _, item := range self.updateQueue {
		uiBase := item.UpdateIF.GetUIBase()
		uiBase.UpdateCount++
	}

	// Update後にRemoveを処理する。Removeがtrueの要素は子供もろとも削除する
	self.root.recRemove()

	return errorList
}

// **********************************************************************
// 再帰実行
// ==================================================
// UpdateTreeの再帰
func (self *UIBase) recUpdateTree(lib *YAMLUI, z int) []error {
	errorList := []error{}

	// 子供の更新
	for _, child := range self.children {
		if child.IsEnable {
			// 更新インターフェースがあれば更新キューに入れる
			if child.updateIF != nil {
				lib.updateQueue = append(lib.updateQueue, UpdateQueueItem{
					UpdateIF: child.updateIF,
					z:        z,
				})
			}

			// UpdateTreeIFがあればそちらを呼び出す。なければ再帰的に呼び出す
			if child.updateTreeIF != nil {
				if err := child.updateTreeIF.UpdateTree(lib, z+1); err != nil {
					errorList = append(errorList, err...)
				}
			} else {
				if err := child.recUpdateTree(lib, z+1); err != nil {
					errorList = append(errorList, err...)
				}
			}
		}
	}

	return errorList
}

// ==================================================
// Removeの再帰
func (self *UIBase) recRemove() {
	for _, child := range self.children {
		if child.Remove {
			// Removeがtrueの要素は子供もろとも削除する
			self.RemoveChild(child)
		} else {
			// 再帰的に呼び出す
			child.recRemove()
		}
	}
}
