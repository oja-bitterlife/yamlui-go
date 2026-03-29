package yamlui

// **********************************************************************
// 左上を (X, Y) とし、幅 W、高さ H の矩形を表す構造体
type Area struct {
	X, Y int
	W, H int
}

// 左端、上端、右端、下端を取得する
func (a Area) Left() int   { return a.X }
func (a Area) Top() int    { return a.Y }
func (a Area) Right() int  { return a.X + a.W }
func (a Area) Bottom() int { return a.Y + a.H }

func NewArea(x, y, w, h int) Area {
	return Area{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
}

// UIBase から Area を取得する
func (ui *UIBase) Area() Area {
	return Area{
		X: ui.X,
		Y: ui.Y,
		W: ui.W,
		H: ui.H,
	}
}

// ==================================================
// ユーティリティー
// 二つのエリアを重ね合わせて、重なった部分のエリアを返す
func (a Area) Clip(limiter Area) Area {
	// 左端：二つの左端のうち、より右にある方
	newX := max(a.Left(), limiter.Left())
	// 上端：二つの上端のうち、より下にある方
	newY := max(a.Top(), limiter.Top())

	// 右端：二つの右端のうち、より左にある方
	newRight := min(a.Right(), limiter.Right())
	// 下端：二つの下端のうち、より上にある方
	newBottom := min(a.Bottom(), limiter.Bottom())

	// 幅と高さを計算（マイナスになったら重なりなしなので 0 にする）
	newW := max(0, newRight-newX)
	newH := max(0, newBottom-newY)

	return Area{
		X: newX,
		Y: newY,
		W: newW,
		H: newH,
	}
}

// 上下左右を指定された量だけ削る（マイナスなら広げる）
func (a Area) Inset(dx, dy int) Area {
	return Area{
		X: a.X + dx,
		Y: a.Y + dy,
		W: max(0, a.W-dx*2),
		H: max(0, a.H-dy*2),
	}
}

// エリアを dx, dy だけ移動させる
func (a Area) Offset(dx, dy int) Area {
	return Area{
		X: a.X + dx,
		Y: a.Y + dy,
		W: a.W,
		H: a.H,
	}
}

// ==================================================
// Getter
// XYを取り出す
func (a Area) XY() (x, y int) {
	return a.X, a.Y
}

// WHを取り出す
func (a Area) WH() (w, h int) {
	return a.W, a.H
}

// XYWHを取り出す
func (a Area) Rect() (x, y, w, h int) {
	return a.X, a.Y, a.W, a.H
}
