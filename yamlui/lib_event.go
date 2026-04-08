package yamlui

import (
	"path"
	"time"
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

// デバッグ用：3秒経っても反応がなければパニックさせる
func (lib *YAMLUI) CallEvent(name string) {
	done := make(chan struct{})
	lib.eventChannel <- eventContext{name: name, done: done}

	select {
	case <-done:
		// 正常終了
	case <-time.After(3 * time.Second):
		panic("CallEvent Deadlock detected! Event: " + name)
	}
}

// **********************************************************************
// イベント処理
// UpdateQueueのEventsにイベントを振り分ける
func (lib *YAMLUI) ProcessEvents(event string) UpdateIF {
	// Updateの後ろからイベント処理
	for i := len(lib.updateQueue) - 1; i >= 0; i-- {
		uiBase := lib.updateQueue[i].UpdateIF.GetUIBase()

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
			return lib.updateQueue[i].UpdateIF
		}
	}
	return nil
}
