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
// Update
type UpdateIF interface {
	Update(lib *YAMLUI, event string) (string, error)
	GetUIBase() *UIBase
}

type UpdateTreeIF interface {
	UpdateTree(lib *YAMLUI, z int) error
}

func (self *UIBase) SetUpdateIF(updateIF UpdateIF) {
	self.updateIF = updateIF
}

func (self *UIBase) SetUpdateTreeIF(updateTreeIF UpdateTreeIF) {
	self.updateTreeIF = updateTreeIF
}

// **********************************************************************
// 呼び出し口
func (lib *YAMLUI) Dispatch(event string) error {
	lib.mtx.Lock()
	defer lib.mtx.Unlock()

	// updateQueueをクリア
	lib.updateQueue = lib.updateQueue[:0]

	// 更新コンテキストを作成してUpdateTreeを呼び出す
	// ----------------------------------------
	if err := lib.root.recUpdateTree(lib, 0); err != nil {
		return err
	}

	// updateQueueに溜まったUpdateを実行する
	// ----------------------------------------
	// Z順でソートする
	slices.SortStableFunc(lib.updateQueue, func(a, b UpdateQueueItem) int {
		return a.z - b.z
	})

	// イベントをUpdateの前に処理し、Updateにイベントを通知する
	updateIF := lib.ProcessEvents(event)
	if updateIF == nil {
		// 処理するUIがなかった
		return nil
	}
	uiBase := updateIF.GetUIBase()

	// スクリプトがあれば走らせる
	// ----------------------------------------
	if uiBase.HasScript() {
		// スクリプトを実行する前に、UIBaseのプロパティをVMに保存しておく
		uiBase.storeToVM(event)

		// スクリプトを実行
		if err := uiBase.script.Run(); err != nil {
			return err
		}

		// スクリプトを実行した後に、VMからUIBaseのプロパティを更新する
		uiBase.loadFromVM()
	}

	// Update
	// ----------------------------------------
	// Updateを呼び出す
	if nextEvent, err := updateIF.Update(lib, event); err != nil {
		return err
	} else if nextEvent != "" {
		// Updateからイベントが返ってきたら処理する
	}

	// Update後にRemoveを処理する。Removeがtrueの要素は子供もろとも削除する
	lib.root.recRemove()

	return nil
}

// **********************************************************************
// 再帰実行
// ==================================================
// UpdateTreeの再帰
func (self *UIBase) recUpdateTree(lib *YAMLUI, z int) error {
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
					return err
				}
			} else {
				if err := child.recUpdateTree(lib, z+1); err != nil {
					return err
				}
			}
		}
	}

	return nil
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
