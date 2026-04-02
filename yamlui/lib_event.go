package yamlui

import (
	"path"
)

// **********************************************************************
// イベント管理
func (self *YAMLUI) AddEvent(event string) {
	self.EventQueue = append(self.EventQueue, event)
}

func (self *YAMLUI) RemoveEvent(event string) {
	for i, e := range self.EventQueue {
		if e == event {
			self.EventQueue = append(self.EventQueue[:i], self.EventQueue[i+1:]...)
			return
		}
	}
}

func (self *YAMLUI) ClearEvents() {
	self.EventQueue = []string{}
}

// **********************************************************************
// イベント処理
// UpdateQueueのEventsにイベントを振り分ける
func (self *YAMLUI) ProcessEvents() {
	// Updateの後ろからイベント処理が可能かを確認する
	checkEvents := self.EventQueue[:]
	for i := len(self.updateQueue) - 1; i >= 0; i-- {
		self.updateQueue[i].ctx.Events = []string{} // イベントリストを初期化

		remainEvents := []string{} // マッチしなかったイベントをためておくリスト
		for _, event := range checkEvents {
			// 受信設定と一致したイベントがあるか確認する
			matchedAny := false
			for _, wildStr := range self.updateQueue[i].ctx.Base.Events {
				// ワイルドカード(path.Match)でマッチング
				if match, _ := path.Match(wildStr, event); match {
					matchedAny = true
					break // 1つでもマッチしたらこのUIのもの！
				}
			}

			// マッチしたイベントはこのUIで処理するためにctx.Eventsに追加する
			if matchedAny {
				self.updateQueue[i].ctx.Events = append(self.updateQueue[i].ctx.Events, event)
			} else {
				// マッチしなかったイベントは次のUIで確認するために残しておく
				remainEvents = append(remainEvents, event)
			}
		}
		checkEvents = remainEvents // 次のUIではマッチしなかったイベントだけを確認する
	}
	self.EventQueue = checkEvents // 最後までマッチしなかったイベントは残しておく
}
