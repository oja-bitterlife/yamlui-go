package yamlui

import (
	"path"
)

// **********************************************************************
// イベント管理
func (self *YAMLUI) AddEvent(event string) {
	self.eventQueue = append(self.eventQueue, event)
}

func (self *YAMLUI) RemoveEvent(event string) {
	for i, e := range self.eventQueue {
		if e == event {
			self.eventQueue = append(self.eventQueue[:i], self.eventQueue[i+1:]...)
			return
		}
	}
}

// ワイルドカードでイベントチェックする便利関数
func (self *YAMLUI) HasEvent(matchStr string, events []string) bool {
	for _, event := range events {
		if match := self.MatchEvent(matchStr, event); match {
			return true
		}
	}
	return false
}

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
	for _, event := range self.eventQueue {

		// Updateの後ろからイベント処理
		for i := len(self.updateQueue) - 1; i >= 0; i-- {

			// 受信設定と一致したイベントがあるか確認する
			matchedAny := false
			for _, wildStr := range self.updateQueue[i].UpdateIF.GetUIBase().Events {

				// ワイルドカード(path.Match)でマッチング
				if match, _ := path.Match(wildStr, event); match {
					matchedAny = true
					break // 1つでもマッチしたらこのUIのもの！
				}
			}

			// マッチしたイベントはupdateQueueのEventsに保存する
			if matchedAny {
				self.updateQueue[i].recv = append(self.updateQueue[i].recv, event)
			}
		}
	}
}
