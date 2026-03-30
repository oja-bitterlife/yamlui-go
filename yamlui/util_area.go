package yamlui

// **********************************************************************
// 左上を (X, Y) とし、幅 W、高さ H の矩形を表す構造体
type Area struct {
	X, Y float64
	W, H float64
}

// 左端、上端、右端、下端を取得する
func (self Area) Left() float64   { return self.X }
func (self Area) Top() float64    { return self.Y }
func (self Area) Right() float64  { return self.X + self.W }
func (self Area) Bottom() float64 { return self.Y + self.H }

func (self Area) ILeft() int   { return int(self.X) }
func (self Area) ITop() int    { return int(self.Y) }
func (self Area) IRight() int  { return int(self.X + self.W) }
func (self Area) IBottom() int { return int(self.Y + self.H) }

func NewArea(x, y, w, h float64) Area {
	return Area{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
}
func NewAreaI(x, y, w, h int) Area {
	return Area{
		X: float64(x),
		Y: float64(y),
		W: float64(w),
		H: float64(h),
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

func (self *UIBase) SetXY(x, y float64) {
	self.X = x
	self.Y = y
}

func (self *UIBase) SetIXY(x, y int) {
	self.X = float64(x)
	self.Y = float64(y)
}

func (self *UIBase) SetWH(w, h float64) {
	self.W = w
	self.H = h
}

func (self *UIBase) SetIWH(w, h int) {
	self.W = float64(w)
	self.H = float64(h)
}

func (self *UIBase) SetRect(x, y, w, h float64) {
	self.SetXY(x, y)
	self.SetWH(w, h)
}

func (self *UIBase) SetIRect(x, y, w, h int) {
	self.SetIXY(x, y)
	self.SetIWH(w, h)
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
func (self Area) Inset(dx, dy float64) Area {
	return Area{
		X: self.X + dx,
		Y: self.Y + dy,
		W: max(0, self.W-dx*2),
		H: max(0, self.H-dy*2),
	}
}

func (self Area) IInset(dx, dy int) Area {
	return Area{
		X: self.X + float64(dx),
		Y: self.Y + float64(dy),
		W: max(0, self.W-float64(dx)*2),
		H: max(0, self.H-float64(dy)*2),
	}
}

// エリアを dx, dy だけ移動させる
func (self Area) Offset(dx, dy float64) Area {
	return Area{
		X: self.X + dx,
		Y: self.Y + dy,
		W: self.W,
		H: self.H,
	}
}

func (self Area) IOffset(dx, dy int) Area {
	return Area{
		X: self.X + float64(dx),
		Y: self.Y + float64(dy),
		W: self.W,
		H: self.H,
	}
}

// ==================================================
// alignment
func (self Area) AlignCenterX(w float64) float64 {
	return self.X + (self.W-w)/2
}

func (self Area) AlignCenterIX(w int) float64 {
	return self.X + (self.W-float64(w))/2
}

func (self Area) AlignCenterY(h float64) float64 {
	return self.Y + (self.H-h)/2
}

func (self Area) AlignCenterIY(h int) float64 {
	return self.Y + (self.H-float64(h))/2
}

func (self Area) AlignCenter(w, h float64) (x, y float64) {
	return self.AlignCenterX(w), self.AlignCenterY(h)
}

func (self Area) AlignCenterI(w, h int) (x, y float64) {
	return self.AlignCenterIX(w), self.AlignCenterIY(h)
}

func (self Area) AlignRight(w float64) float64 {
	return self.X + self.W - w
}

func (self Area) AlignIRight(w int) float64 {
	return self.X + self.W - float64(w)
}

func (self Area) AlignBottom(h float64) float64 {
	return self.Y + self.H - h
}

func (self Area) AlignIBottom(h int) float64 {
	return self.Y + self.H - float64(h)
}

func (self Area) AlignRightBottom(w, h float64) (x, y float64) {
	return self.AlignRight(w), self.AlignBottom(h)
}

func (self Area) AlignIRightBottom(w, h int) (x, y float64) {
	return self.AlignIRight(w), self.AlignIBottom(h)
}

// ==================================================
// Getter
// int型でX、Yを取り出す
func (self Area) IX() int {
	return int(self.X)
}

func (self Area) IY() int {
	return int(self.Y)
}

// XYを取り出す
func (self Area) XY() (x, y float64) {
	return self.X, self.Y
}

func (self Area) IXY() (x, y int) {
	return int(self.X), int(self.Y)
}

// WHを取り出す
func (self Area) WH() (w, h float64) {
	return self.W, self.H
}

func (self Area) IWH() (w, h int) {
	return int(self.W), int(self.H)
}

// XYWHを取り出す
func (self Area) Rect() (x, y, w, h float64) {
	return self.X, self.Y, self.W, self.H
}

func (self Area) IRect() (x, y, w, h int) {
	return int(self.X), int(self.Y), int(self.W), int(self.H)
}

// ==================================================
// ユーティリティー関数
func (self Area) Contains(x, y float64) bool {
	return x >= self.Left() && x < self.Right() && y >= self.Top() && y < self.Bottom()
}

func (self Area) IContains(x, y int) bool {
	return self.Contains(float64(x), float64(y))
}

func (self Area) IsEmpty() bool {
	return self.W <= 0 || self.H <= 0
}
