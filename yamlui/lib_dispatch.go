package yamlui

import (
	"slices"
)

// **********************************************************************
// Dispatchのインターフェース
// ==================================================
// 直接DispatchIFを呼び出すのではなく、DispatchTreeでDispatchQueueItemにしてキューに入れる
type DispatchQueueItem struct {
	DispatchIF DispatchIF
	z          int
}

// ==================================================
// Update
type DispatchIF interface {
	Dispatch(lib *YAMLUI, event string) (string, error)
	GetUIBase() *UIBase
}

type DispatchTreeIF interface {
	DispatchTree(lib *YAMLUI, z int) error
}

func (self *UIBase) SetDispatchIF(dispatchIF DispatchIF) {
	self.dispatchIF = dispatchIF
}

func (self *UIBase) SetDispatchTreeIF(dispatchTreeIF DispatchTreeIF) {
	self.dispatchTreeIF = dispatchTreeIF
}

// **********************************************************************
// 呼び出し口
func (lib *YAMLUI) dispatch(event string) error {
	// 最初にLock
	lib.mtx.Lock()
	lib.isLock.Store(true)
	defer func() {
		lib.mtx.Unlock()
		lib.isLock.Store(false)
	}()

	// queueをクリア
	lib.dispatchQueue = lib.dispatchQueue[:0]

	// 更新コンテキストを作成してDispatchTreeを呼び出す
	// ----------------------------------------
	if err := lib.root.recDispatchTree(lib, 0); err != nil {
		return err
	}

	// queueに溜まったDispatchを実行する
	// ----------------------------------------
	// Z順でソートする
	slices.SortStableFunc(lib.dispatchQueue, func(a, b DispatchQueueItem) int {
		return a.z - b.z
	})

	// イベントをDispatchの前に処理しイベントを通知する対象を決定する
	dispatchIF := lib.ProcessEvents(event)
	if dispatchIF == nil {
		// 処理するUIがなかった
		return nil
	}
	uiBase := dispatchIF.GetUIBase()

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
	// Dispatchを呼び出す
	if nextEvent, err := dispatchIF.Dispatch(lib, event); err != nil {
		return err
	} else if nextEvent != "" {
		lib.SendEvent(nextEvent)
	}

	// Update後にRemoveを処理する。Removeがtrueの要素は子供もろとも削除する
	lib.root.recRemove()

	return nil
}

// **********************************************************************
// 再帰実行
// ==================================================
// UpdateTreeの再帰
func (self *UIBase) recDispatchTree(lib *YAMLUI, z int) error {
	// 子供の更新
	for _, child := range self.children {
		if child.IsEnable {
			// 更新インターフェースがあれば更新キューに入れる
			if child.dispatchIF != nil {
				lib.dispatchQueue = append(lib.dispatchQueue, DispatchQueueItem{
					DispatchIF: child.dispatchIF,
					z:          z,
				})
			}

			// UpdateTreeIFがあればそちらを呼び出す。なければ再帰的に呼び出す
			if child.dispatchTreeIF != nil {
				if err := child.dispatchTreeIF.DispatchTree(lib, z+1); err != nil {
					return err
				}
			} else {
				if err := child.recDispatchTree(lib, z+1); err != nil {
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
