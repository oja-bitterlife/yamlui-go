package yamlui

import (
	"path"
)

func (lib *YAMLUI) SendEvent(event string) {
	// goroutineでイベントを送る
	lib.eventChannel <- event
}

// ==================================================
// イベントチェック.ワイルドカードでイベントチェックする
// 発生したイベントの中に、matchStrにマッチするものがあるか
func (lib *YAMLUI) HasEvent(matchStr string, events []string) bool {
	for _, event := range events {
		if match := lib.MatchEvent(matchStr, event); match {
			return true
		}
	}
	return false
}

// そのイベントがワイルドカードにマッチするか
func (lib *YAMLUI) MatchEvent(matchStr string, event string) bool {
	if match, _ := path.Match(matchStr, event); match {
		return true
	}
	return false
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
