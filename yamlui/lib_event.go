package yamlui

import (
	"path"
)

// **********************************************************************
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
