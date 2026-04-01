package yamlui

type UIEventQueue struct {
	Events []string
}

func (self *UIEventQueue) Add(event string) {
	self.Events = append(self.Events, event)
}

func (self *UIEventQueue) Remove(event string) {
	for i, e := range self.Events {
		if e == event {
			self.Events = append(self.Events[:i], self.Events[i+1:]...)
			return
		}
	}
}
