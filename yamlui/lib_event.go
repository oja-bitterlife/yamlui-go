package yamlui

import (
	"path"
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
// 再帰でTreeをDispatchQueueに入れる
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

// **********************************************************************
// イベント処理
// Dispatch対象のUIを探す
func (lib *YAMLUI) ProcessEvents(event string) DispatchIF {
	// Updateの後ろからイベント処理
	for i := len(lib.dispatchQueue) - 1; i >= 0; i-- {
		uiBase := lib.dispatchQueue[i].DispatchIF.GetUIBase()

		// 受信設定と一致したイベントがあるか確認する
		matchedAny := false
		for _, wildStr := range uiBase.Events {

			// ワイルドカード(path.Match)でマッチング
			if match, _ := path.Match(wildStr, event); match {
				matchedAny = true
				break // 1つでもマッチしたらこのUIのもの！
			}
		}

		// マッチしたイベントはUpdateQueueのどれに対応するかセットする
		if matchedAny {
			return lib.dispatchQueue[i].DispatchIF
		}
	}
	return nil
}

// **********************************************************************
// goroutine用
// ==================================================
// 完了検知のための構造体
type eventContext struct {
	name string
	done chan struct{}
}

// イベント送信
func (lib *YAMLUI) SendEvent(event string) {
	// goroutineでイベントを送る
	lib.eventChannel <- eventContext{
		name: event,
		done: nil,
	}
}

// イベントを送信して完了を待つ
func (lib *YAMLUI) CallEvent(name string) {
	if lib.isLock.Load() {
		panic("Updaate/Draw中にイベントを呼び出すことはできません")
	}

	done := make(chan struct{})
	lib.eventChannel <- eventContext{name: name, done: done}
	<-done
}

// そもそもgoroutinを使わない
func (lib *YAMLUI) DispatchEvent(name string) {
	if lib.isLock.Load() {
		panic("Updaate/Draw中にイベントを呼び出すことはできません")
	}
	lib.dispatch(name)
}
