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
	uiBase     *UIBase
}

// **********************************************************************
// 呼び出し口
func (lib *YAMLUI) dispatch(event string) {
	// 最初にLock
	lib.mtx.Lock()
	lib.isLock.Store(true)
	defer func() {
		lib.mtx.Unlock()
		lib.isLock.Store(false)
	}()

	// DispatchTreeを呼び出してDispatchQueueに入れる
	// ----------------------------------------
	lib.dispatchQueue = lib.dispatchQueue[:0] // queueをクリア
	lib.root.recDispatchTree(lib, 0)

	// queueに溜まったDispatchを実行する
	// ----------------------------------------
	// Z順でソートする.兄弟が先に呼ばれるようになる
	slices.SortStableFunc(lib.dispatchQueue, func(a, b DispatchQueueItem) int {
		return a.z - b.z
	})

	// イベントをDispatchの前に処理しイベントを通知する対象を決定する
	dispatchIF, uiBase := lib.processEvents(event)
	if dispatchIF == nil || uiBase == nil {
		// 処理するUIがなかった
		return
	}

	// スクリプトがあれば走らせる
	// ----------------------------------------
	if uiBase.HasScript() {
		// スクリプトを実行する前に、いくつかのプロパティをVMに保存しておく
		uiBase.storeToVM(lib.Frame, event)

		// スクリプトを実行
		uiBase.script.Run()
	}

	// Update
	// ----------------------------------------
	// Dispatchを呼び出す
	dispatchIF.Dispatch(lib, event)

	// Update後にRemoveを処理する。Removeがtrueの要素は子供もろとも削除する
	lib.root.recRemove()
}

// **********************************************************************
// 再帰でTreeをDispatchQueueに入れる
func (self *UIBase) recDispatchTree(lib *YAMLUI, z int) {
	// 子供の更新
	for _, child := range self.children {
		if child.PropBool(PROP_IS_ENABLE) {
			// 更新インターフェースがあれば更新キューに入れる
			if child.dispatchIF != nil {
				lib.dispatchQueue = append(lib.dispatchQueue, DispatchQueueItem{
					DispatchIF: child.dispatchIF,
					z:          z,
					uiBase:     child,
				})
			}

			// 再帰的に呼び出す
			child.recDispatchTree(lib, z+1)
		}
	}
}

// ==================================================
// Removeの再帰
func (self *UIBase) recRemove() {
	for _, child := range self.children {
		if child.PropBool(PROP_REMOVE) {
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
func (lib *YAMLUI) MatchEvent(wildStr string, event string) bool {
	match, err := path.Match(wildStr, event)
	if err != nil {
		lib.LogErr("Invalid event pattern: %s", wildStr)
	}
	return match
}

// Dispatch対象のUIを探す
func (lib *YAMLUI) processEvents(event string) (DispatchIF, *UIBase) {
	// Updateの後ろからイベント処理
	for i := len(lib.dispatchQueue) - 1; i >= 0; i-- {
		uiBase := lib.dispatchQueue[i].uiBase

		// 受信設定と一致したイベントがあるか確認する
		for _, wildStr := range uiBase.Events {
			// ワイルドカード(path.Match)でマッチング
			if match := lib.MatchEvent(wildStr, event); match {
				// 1つでもマッチしたらこのUIのもの！
				return lib.dispatchQueue[i].DispatchIF, uiBase
			}
		}
	}
	return nil, nil
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
	// Startしていない
	if !lib.isRunning.Load() {
		lib.LogErr("YAMLUI is not running. Call Start() before sending events.")
		return
	}

	// goroutineでイベントを送る
	lib.eventChannel <- eventContext{
		name: event,
		done: nil,
	}
}

// イベントを送信して完了を待つ
func (lib *YAMLUI) CallEvent(event string) {
	// Startしていない
	if !lib.isRunning.Load() {
		// UIのgoruoutineが走っていないときは直接呼び出す
		lib.LogErr("YAMLUI is not running. Call Start() before calling events. Calling event directly.")
		lib.DispatchEvent(event)
		return
	}

	if lib.isLock.Load() {
		panic("Updaate/Draw中にイベントを呼び出すことはできません")
	}

	done := make(chan struct{})
	lib.eventChannel <- eventContext{name: event, done: done}
	<-done
}

// そもそもgoroutinを使わない
func (lib *YAMLUI) DispatchEvent(event string) {
	if lib.isLock.Load() {
		panic("Updaate/Draw中にイベントを呼び出すことはできません")
	}
	lib.dispatch(event)
}
