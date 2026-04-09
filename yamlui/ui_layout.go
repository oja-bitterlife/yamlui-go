package yamlui

import (
	"github.com/oja-bitterlife/yamlui-go/script"
)

type UILayout struct {
	UIBase *UIBase

	// Layout用のプロパティ.ここに置くとprivate扱いになる
	// publicにしたい場合はGetUIBase().SetPropで保存する
	// ==================================================
	// Margin用
	Margin       int
	MarginTop    int
	MarginBottom int
	MarginLeft   int
	MarginRight  int
	MarginX      int
	MarginY      int

	// Align系のプロパティ
	AlignCenter  bool
	AlignCenterX bool
	AlignCenterY bool
	AlignRight   bool
	AlignBottom  bool

	// 自分のXYは絶対座標
	IsAbs bool
}

func NewUILayout(lib *YAMLUI) *UILayout {
	return &UILayout{
		UIBase: NewUIBase(lib),
	}
}

// **********************************************************************
// UIComponentIFの実装
func (self *UILayout) GetUIBase() *UIBase {
	return self.UIBase
}

func (self *UILayout) Clone() UICloned {
	newLayout := *self
	newLayout.UIBase = self.UIBase.Clone()
	return &newLayout
}

func (self *UILayout) Setup(type_ string, data script.ValueMap) error {
	if err := self.UIBase.Setup(type_, data); err != nil { // super call
		return err
	}

	// DrawTreeIFをセット
	self.UIBase.SetDrawTreeFilterIF(self)

	// デフォルトはMargin=0
	self.Margin = data.GetNum("Margin", 0)
	self.MarginTop = data.GetNum("MarginTop", 0)
	self.MarginBottom = data.GetNum("MarginBottom", 0)
	self.MarginLeft = data.GetNum("MarginLeft", 0)
	self.MarginRight = data.GetNum("MarginRight", 0)
	self.MarginX = data.GetNum("MarginX", 0)
	self.MarginY = data.GetNum("MarginY", 0)

	// Align系のプロパティはデフォルトはfalse
	self.AlignCenter = data.GetBool("AlignCenter", false)
	self.AlignCenterX = data.GetBool("AlignCenterX", false)
	self.AlignCenterY = data.GetBool("AlignCenterY", false)
	self.AlignRight = data.GetBool("AlignRight", false)
	self.AlignBottom = data.GetBool("AlignBottom", false)

	// IsAbsのデフォルトはfalse
	self.IsAbs = data.GetBool("IsAbs", false)

	return nil
}

// **********************************************************************
// DrawTreeのOverride
func (self *UILayout) DrawTreeFilter(lib *YAMLUI, ctx DrawTreeContext) DrawTreeContext {
	// 面倒なので先にintにしておく
	parentW, parentH := ctx.ParentClip.WH()
	selfW, selfH := self.GetUIBase().Area().WH()

	// 絶対座標対応
	if self.IsAbs {
		ctx.X = lib.Screen.X
		ctx.Y = lib.Screen.Y
	}

	// Align系のプロパティを考慮して、子の描画領域を計算する
	offsetX := 0
	offsetY := 0
	if self.AlignRight {
		offsetX = parentW - selfW
	}
	if self.AlignBottom {
		offsetY = parentH - selfH
	}
	if self.AlignCenter || self.AlignCenterX {
		offsetX = (parentW - selfW) / 2
	}
	if self.AlignCenter || self.AlignCenterY {
		offsetY = (parentH - selfH) / 2
	}
	// alignされた座標
	alignX := ctx.X + offsetX
	alignY := ctx.Y + offsetY

	// マージンを考慮して、子の描画領域を計算する
	left := self.MarginLeft + self.MarginX + self.Margin
	top := self.MarginTop + self.MarginY + self.Margin
	right := self.MarginRight + self.MarginX + self.Margin
	bottom := self.MarginBottom + self.MarginY + self.Margin

	// Right/Bottomマージン分親のエリアから引いておく
	parentArea := ctx.ParentClip
	parentArea.W -= right
	parentArea.H -= bottom

	// Left/Topマージン分座標をずらす
	newClip := Area{
		X: alignX + left,
		Y: alignY + top,
		W: selfW,
		H: selfH,
	}.Clip(parentArea)

	// コンテキストの更新
	ctx.X = alignX + left
	ctx.Y = alignY + top
	ctx.Clip = newClip

	// 更新したコンテキストを返す
	return ctx
}
