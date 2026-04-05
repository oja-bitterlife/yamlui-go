package yamlui

// **********************************************************************
// 左上を (X, Y) とし、幅 W、高さ H の矩形を表す構造体
type Area struct {
	X, Y int
	W, H int
}

// 左端、上端、右端、下端を取得する
func (self Area) Left() int   { return self.X }
func (self Area) Top() int    { return self.Y }
func (self Area) Right() int  { return self.X + self.W }
func (self Area) Bottom() int { return self.Y + self.H }

func NewArea(x, y, w, h int) Area {
	return Area{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
}

// ==================================================
// UIBase用機能
// UIBase から Area を取得する
func (self *UIBase) Area() Area {
	return Area{
		X: self.X,
		Y: self.Y,
		W: self.W,
		H: self.H,
	}
}

func (self *UIBase) SetXY(x, y int) {
	self.X = x
	self.Y = y
}

func (self *UIBase) SetWH(w, h int) {
	self.W = w
	self.H = h
}

func (self *UIBase) SetRect(x, y, w, h int) {
	self.SetXY(x, y)
	self.SetWH(w, h)
}

// ==================================================
// ユーティリティー
// 二つのエリアを重ね合わせて、重なった部分のエリアを返す
func (self Area) Clip(limiter Area) Area {
	// 左端：二つの左端のうち、より右にある方
	newX := max(self.Left(), limiter.Left())
	// 上端：二つの上端のうち、より下にある方
	newY := max(self.Top(), limiter.Top())

	// 右端：二つの右端のうち、より左にある方
	newRight := min(self.Right(), limiter.Right())
	// 下端：二つの下端のうち、より上にある方
	newBottom := min(self.Bottom(), limiter.Bottom())

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
func (self Area) Inset(dx, dy int) Area {
	return Area{
		X: self.X + dx,
		Y: self.Y + dy,
		W: max(0, self.W-dx*2),
		H: max(0, self.H-dy*2),
	}
}

// エリアを dx, dy だけ移動させる
func (self Area) Offset(dx, dy int) Area {
	return Area{
		X: self.X + dx,
		Y: self.Y + dy,
		W: self.W,
		H: self.H,
	}
}

// ==================================================
// alignment
func (self Area) AlignCenterX(w int) int {
	return self.X + (self.W-w)/2
}

func (self Area) AlignCenterY(h int) int {
	return self.Y + (self.H-h)/2
}

func (self Area) AlignCenter(w, h int) (x, y int) {
	return self.AlignCenterX(w), self.AlignCenterY(h)
}

func (self Area) AlignRight(w int) int {
	return self.X + self.W - w
}

func (self Area) AlignBottom(h int) int {
	return self.Y + self.H - h
}

func (self Area) AlignRightBottom(w, h int) (x, y int) {
	return self.AlignRight(w), self.AlignBottom(h)
}

// ==================================================
// Getter
// XYを取り出す
func (self Area) XY() (x, y int) {
	return self.X, self.Y
}

// WHを取り出す
func (self Area) WH() (w, h int) {
	return self.W, self.H
}

// XYWHを取り出す
func (self Area) Rect() (x, y, w, h int) {
	return self.X, self.Y, self.W, self.H
}

// ==================================================
// ユーティリティー関数
func (self Area) Contains(x, y int) bool {
	return x >= self.Left() && x < self.Right() && y >= self.Top() && y < self.Bottom()
}

func (self Area) ContainsX(x int) bool {
	return x >= self.Left() && x < self.Right()
}

func (self Area) ContainsY(y int) bool {
	return y >= self.Top() && y < self.Bottom()
}

func (self Area) IsEmpty() bool {
	return self.W <= 0 || self.H <= 0
}
