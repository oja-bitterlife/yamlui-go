package yamlui

import (
	"path"
)

// **********************************************************************
// イベント管理
// ==================================================
// イベントキュー
type EventQueueItem struct {
	event         string
	updateQueueID string
}

func (lib *YAMLUI) AddEvent(event string) {
	lib.eventQueue = append(lib.eventQueue, EventQueueItem{
		event:         event,
		updateQueueID: "", // UpdateQueueIDはProcessEventsでセットされる
	})
}

func (lib *YAMLUI) ReserveEvent(event string) {
	lib.eventReserve = append(lib.eventReserve, EventQueueItem{
		event:         event,
		updateQueueID: "", // UpdateQueueIDはProcessEventsでセットされる
	})
}

func (lib *YAMLUI) RemoveEvent(event string) {
	for i, e := range lib.eventQueue {
		if e.event == event {
			lib.eventQueue = append(lib.eventQueue[:i], lib.eventQueue[i+1:]...)
			return
		}
	}
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
func (lib *YAMLUI) ProcessEvents() {
	// イベントを順番にUpdateQueueの後ろから振り分ける
	for ei, e := range lib.eventQueue {

		// Updateの後ろからイベント処理
		for i := len(lib.updateQueue) - 1; i >= 0; i-- {
			uiBase := lib.updateQueue[i].UpdateIF.GetUIBase()

			// 受信設定と一致したイベントがあるか確認する
			matchedAny := false
			for _, wildStr := range uiBase.Events {

				// ワイルドカード(path.Match)でマッチング
				if match, _ := path.Match(wildStr, e.event); match {
					matchedAny = true
					break // 1つでもマッチしたらこのUIのもの！
				}
			}

			// マッチしたイベントはUpdateQueueのどれに対応するかセットする
			if matchedAny {
				lib.eventQueue[ei].updateQueueID = uiBase.ID
			}
		}
	}
}
