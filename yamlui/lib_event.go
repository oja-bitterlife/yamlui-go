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

func (self *YAMLUI) ClearEvents() {
	self.eventQueue = []string{}
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
	// Updateの後ろからイベント処理が可能かを確認する
	checkEvents := self.eventQueue[:]
	for i := len(self.updateQueue) - 1; i >= 0; i-- {
		remainEvents := []string{} // マッチしなかったイベントをためておくリスト
		for _, event := range checkEvents {
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
				self.updateQueue[i].Events = append(self.updateQueue[i].Events, event)
			} else {
				// マッチしなかったイベントは次のUIで確認するために残しておく
				remainEvents = append(remainEvents, event)
			}
		}
		checkEvents = remainEvents // 次のUIではマッチしなかったイベントだけを確認する
	}
	self.eventQueue = checkEvents // 最後までマッチしなかったイベントは残しておく
}
