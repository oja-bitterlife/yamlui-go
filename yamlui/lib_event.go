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
	updateQueueNo int // 初期値は-1。ProcessEventsでUpdateQueueのどれに対応するかセットされる
}

func (self *YAMLUI) AddEvent(event string) {
	self.eventQueue = append(self.eventQueue, EventQueueItem{
		event:         event,
		updateQueueNo: -1, // UpdateQueueNoはProcessEventsでセットされる
	})
}

func (self *YAMLUI) RemoveEvent(event string) {
	for i, e := range self.eventQueue {
		if e.event == event {
			self.eventQueue = append(self.eventQueue[:i], self.eventQueue[i+1:]...)
			return
		}
	}
}

// ==================================================
// イベントチェック.ワイルドカードでイベントチェックする
// 発生したイベントの中に、matchStrにマッチするものがあるか
func (self *YAMLUI) HasEvent(matchStr string, events []string) bool {
	for _, event := range events {
		if match := self.MatchEvent(matchStr, event); match {
			return true
		}
	}
	return false
}

// そのイベントがワイルドカードにマッチするか
func (self *YAMLUI) MatchEvent(matchStr string, event string) bool {
	if match, _ := path.Match(matchStr, event); match {
		return true
	}
	return false
}

// **********************************************************************
// イベント処理
// UpdateQueueのEventsにイベントを振り分ける
func (self *YAMLUI) ProcessEvents() {
	// イベントを順番にUpdateQueueの後ろから振り分ける
	for ei, e := range self.eventQueue {

		// Updateの後ろからイベント処理
		for i := len(self.updateQueue) - 1; i >= 0; i-- {

			// 受信設定と一致したイベントがあるか確認する
			matchedAny := false
			for _, wildStr := range self.updateQueue[i].UpdateIF.GetUIBase().Events {

				// ワイルドカード(path.Match)でマッチング
				if match, _ := path.Match(wildStr, e.event); match {
					matchedAny = true
					break // 1つでもマッチしたらこのUIのもの！
				}
			}

			// マッチしたイベントはUpdateQueueのどれに対応するかセットする
			if matchedAny {
				self.eventQueue[ei].updateQueueNo = i
			}
		}
	}
}
