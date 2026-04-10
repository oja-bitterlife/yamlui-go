package yamlui

import (
	"slices"
)

// **********************************************************************
// 描画インターフェース
// ==================================================
// 直接DrawIFを呼び出すのではなく、DrawTreeの中でDrawQueueItemにしてキューに入れる
type DrawQueueItem struct {
	drawIF DrawIF
	z      float64
	x, y   int
	clip   Area
}

// ==================================================
// DrawTreeIF
// Drawの引数やZ順を変更したいときに使う
type DrawContext struct {
	Parent     *UIBase
	ParentClip Area
	Z          float64
	X, Y       int
	Clip       Area
}

type DrawContextIF interface {
	DrawContextFilter(lib *YAMLUI, ctx DrawContext) DrawContext
}

// 通常のDraw
type DrawIF interface {
	Draw(lib *YAMLUI, x, y int, clip Area)
}

// 実装確認君
func (self *UIBase) CheckDrawContextIF(drawContextIF DrawContextIF) {}
func (self *UIBase) CheckDrawIF(drawIF DrawIF)                      {}

// **********************************************************************
// 呼び出し口
func (lib *YAMLUI) Draw(sx, sy, sw, sh int) {
	// 最初にLock
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
	ctx := DrawContext{
		Parent:     nil,
		ParentClip: lib.Screen,
		Z:          0,
		X:          sx,
		Y:          sy,
		Clip:       lib.Screen,
	}
	lib.root.recDrawTree(lib, ctx)

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
		item.drawIF.Draw(lib, item.x, item.y, item.clip)
	}

	lib.Frame++
}

// **********************************************************************
// DrawTreeの再帰実行
func (self *UIBase) recDrawTree(lib *YAMLUI, ctx DrawContext) {
	// Filterがあれば描画コンテキストを変更する
	if self.drawTreeIF != nil {
		ctx = self.drawTreeIF.DrawContextFilter(lib, ctx)
	}

	// 描画インターフェースがあれば描画キューに入れる
	if self.drawIF != nil {
		lib.drawQueue = append(lib.drawQueue, DrawQueueItem{
			drawIF: self.drawIF,
			z:      ctx.Z,
			x:      ctx.X,
			y:      ctx.Y,
			clip:   ctx.Clip,
		})
	}

	// 子供の描画
	for _, child := range self.children {
		if child.PropBool(PROP_IS_VISIBLE) {
			childArea := child.Area().Offset(ctx.X, ctx.Y)
			childClip := childArea.Clip(ctx.Clip)

			// 子供の描画コンテキストを作る。Z順は親よりも大きくする
			childCtx := DrawContext{
				Parent:     self,
				ParentClip: ctx.Clip,
				Z:          ctx.Z + 1,
				X:          childArea.X,
				Y:          childArea.Y,
				Clip:       childClip,
			}

			// 再帰的に子供を描画
			child.recDrawTree(lib, childCtx)
		}
	}
}
