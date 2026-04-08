package yamlui

import "slices"

// **********************************************************************
// 描画インターフェース
type DrawContext struct {
	Lib        *YAMLUI // ライブラリ全体へのアクセス
	Parent     *UIBase // 親
	ParentClip Area    // 親の描画領域（クリップされている）
}

// 直接DrawIFを呼び出すのではなく、DrawTreeの中でDrawQueueItemにしてキューに入れる
type DrawQueueItem struct {
	drawIF DrawIF
	z      float64
	x      int
	y      int
	clip   Area
	ctx    DrawContext
}

// 通常のDraw
type DrawIF interface {
	Draw(x, y int, clip Area, ctx DrawContext)
}

// Align操作が必要な時とか、DrawTreeを自前で実装したいときのインターフェース
type DrawTreeIF interface {
	DrawTree(z float64, x, y int, clip Area, ctx DrawContext)
}

func (self *UIBase) SetDrawIF(drawIF DrawIF) {
	self.drawIF = drawIF
}

func (self *UIBase) SetDrawTreeIF(drawTreeIF DrawTreeIF) {
	self.drawTreeIF = drawTreeIF
}

// コンテキスト作成
func NewDrawContext(lib *YAMLUI, parent *UIBase, parentClip Area) DrawContext {
	return DrawContext{
		Lib:        lib,
		Parent:     parent,
		ParentClip: parentClip,
	}
}

// **********************************************************************
// 呼び出し口
func (lib *YAMLUI) Draw(sx, sy, sw, sh int) {
	lib.mtx.RLock()
	lib.isLock.Store(true)
	defer func() {
		lib.mtx.RUnlock()
		lib.isLock.Store(false)
	}()

	// 描画領域をセット
	lib.Screen = NewArea(sx, sy, sw, sh)

	// drawQueueをクリア
	lib.drawQueue = lib.drawQueue[:0]

	// 描画コンテキストを作成してDrawTreeを呼び出す
	// ----------------------------------------
	ctx := NewDrawContext(lib, nil, lib.Screen)
	lib.root.RecDrawTree(0, sx, sy, lib.Screen, ctx)

	// drawQueueに溜まった描画命令を実行する
	// ----------------------------------------
	// Z順でソートする
	slices.SortStableFunc(lib.drawQueue, func(a, b DrawQueueItem) int {
		if a.z < b.z {
			return -1
		}
		if a.z > b.z {
			return 1
		}
		return 0
	})

	// ソートされたqueueを順番に実行する
	for _, item := range lib.drawQueue {
		item.drawIF.Draw(item.x, item.y, item.clip, item.ctx)
	}
}

// **********************************************************************
// DrawTreeの再帰実行
func (self *UIBase) RecDrawTree(z float64, x, y int, clip Area, ctx DrawContext) {
	// 子供の描画
	for _, child := range self.children {
		if child.IsVisible {
			// 描画座標やクリップ領域等を計算してコンテキストを作る
			childCtx := NewDrawContext(ctx.Lib, self, clip)
			childArea := child.Area().Offset(x, y)
			childClip := childArea.Clip(clip)

			// 描画インターフェースがあれば描画キューに入れる
			if child.drawIF != nil {
				ctx.Lib.drawQueue = append(ctx.Lib.drawQueue, DrawQueueItem{
					drawIF: child.drawIF,
					z:      z,
					x:      childArea.X,
					y:      childArea.Y,
					clip:   childClip,
					ctx:    childCtx,
				})
			}

			// 自前drawTreeがあればそちらを呼び出す
			if child.drawTreeIF != nil {
				child.drawTreeIF.DrawTree(z+1, childArea.X, childArea.Y, childClip, childCtx)
			} else {
				// なければ再帰的に子供を描画
				child.RecDrawTree(z+1, childArea.X, childArea.Y, childClip, childCtx)
			}
		}
	}
}
