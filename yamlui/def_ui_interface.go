package yamlui

// **********************************************************************
// Dispatch/Drawのインターフェース

// ==================================================
// Update
// これを実装するとEventをキャッチした時に呼び出される
// UI作成時はEVENT_UI_ONCREATEが呼び出される
type DispatchIF interface {
	Dispatch(lib *YAMLUI, event string)
}

// 実装確認君
func (self *UIBase) CheckDispatchIF(dispatchIF DispatchIF) {}
func (lib *YAMLUI) CheckDispatchIF(dispatchIF DispatchIF)  {}

// ==================================================
// Draw
type DrawContext struct {
	Parent     *UIBase // 親のUIBase
	ParentClip Area    // 親のClip。親のClipは子のClipの上限になる
	Z          float64 // Z順。大きいほど前に描画される
	X, Y       int     // 絶対座標の描画位置。相対座標はPROPに入っている(Area()で取り出す)
	Clip       Area    // 描画領域
}

// Drawの引数やZ順を変更したいときに使う.UITreeのトラバース中に呼び出される
type DrawContextIF interface {
	DrawContextFilter(lib *YAMLUI, ctx DrawContext) DrawContext
}

// これを実装するとDrawのときに呼び出される
type DrawIF interface {
	Draw(lib *YAMLUI, x, y int, clip Area)
}

// 実装確認君
func (self *UIBase) CheckDrawContextIF(drawContextIF DrawContextIF) {}
func (lib *YAMLUI) CheckDrawContextIF(drawContextIF DrawContextIF)  {}

func (self *UIBase) CheckDrawIF(drawIF DrawIF) {}
func (lib *YAMLUI) CheckDrawIF(drawIF DrawIF)  {}
