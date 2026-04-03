package yamlui

// ==================================================
// 更新
// ----------------------------------------
// Updateのインターフェース
type OnInitIF interface {
	OnInit(lib *YAMLUI) error
}

type UpdateIF interface {
	Update(lib *YAMLUI, events []string) error
	GetUIBase() *UIBase
}

type UpdateTreeIF interface {
	UpdateTree(lib *YAMLUI, z int) []error
}

// 直接DrawIFを呼び出すのではなく、UpdateTreeの中でUpdateQueueItemにしてキューに入れる
type UpdateQueueItem struct {
	UpdateIF UpdateIF
	z        int
	Events   []string
}

func (self *UIBase) SetUpdateIF(updateIF UpdateIF) {
	self.updateIF = updateIF
}

func (self *UIBase) SetUpdateTreeIF(updateTreeIF UpdateTreeIF) {
	self.updateTreeIF = updateTreeIF
}

// ==================================================
// Update実行
// 再帰実行
func (self *UIBase) RecUpdateTree(lib *YAMLUI, z int) []error {
	errorList := []error{}

	// 子供の更新
	for _, child := range self.children {
		if child.IsEnable {
			// 更新インターフェースがあれば更新キューに入れる
			if child.updateIF != nil {
				lib.updateQueue = append(lib.updateQueue, UpdateQueueItem{
					UpdateIF: child.updateIF,
					z:        z,
					Events:   []string{}, // イベントはProcessEventsで振り分ける
				})
			}

			// UpdateTreeIFがあればそちらを呼び出す。なければ再帰的に呼び出す
			if child.updateTreeIF != nil {
				if err := child.updateTreeIF.UpdateTree(lib, z+1); err != nil {
					errorList = append(errorList, err...)
				}
			} else {
				if err := child.RecUpdateTree(lib, z+1); err != nil {
					errorList = append(errorList, err...)
				}
			}
		}
	}

	return errorList
}
