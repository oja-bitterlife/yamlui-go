package yamlui

import "slices"

// **********************************************************************
// Updateのインターフェース
type OnInitIF interface {
	OnInit(lib *YAMLUI) error
}

type UpdateIF interface {
	Update(lib *YAMLUI, events []string) error
	GetUIBase() *UIBase
}

type UpdateTreeIF interface {
	UpdateTree(lib *YAMLUI, z int) []error
}

// 直接DrawIFを呼び出すのではなく、UpdateTreeの中でUpdateQueueItemにしてキューに入れる
type UpdateQueueItem struct {
	UpdateIF UpdateIF
	z        int
	Events   []string
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
	self.updateQueue = []UpdateQueueItem{}

	// 更新コンテキストを作成してUpdateTreeを呼び出す
	// ----------------------------------------
	if err := self.root.RecUpdateTree(self, 0); err != nil {
		errorList = append(errorList, err...)
	}

	// updateQueueに溜まった描画命令を実行する
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

	// ここ以降でself.eventQueueに入るイベントはUpateで入るイベント

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
			uiBase.storeScriptEvent(item.Events) // @UIEventを追加

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

// **********************************************************************
// 再帰実行
func (self *UIBase) RecUpdateTree(lib *YAMLUI, z int) []error {
	errorList := []error{}

	// 子供の更新
	for _, child := range self.children {
		if child.IsEnable {
			// 更新インターフェースがあれば更新キューに入れる
			if child.updateIF != nil {
				lib.updateQueue = append(lib.updateQueue, UpdateQueueItem{
					UpdateIF: child.updateIF,
					z:        z,
					Events:   []string{}, // イベントはProcessEventsで振り分ける
				})
			}

			// UpdateTreeIFがあればそちらを呼び出す。なければ再帰的に呼び出す
			if child.updateTreeIF != nil {
				if err := child.updateTreeIF.UpdateTree(lib, z+1); err != nil {
					errorList = append(errorList, err...)
				}
			} else {
				if err := child.RecUpdateTree(lib, z+1); err != nil {
					errorList = append(errorList, err...)
				}
			}
		}
	}

	return errorList
}
