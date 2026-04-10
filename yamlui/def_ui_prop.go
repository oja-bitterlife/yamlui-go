package yamlui

// **********************************************************************
// UIBaseのプロパティ名
// YAMLのプロパティで使えるので一覧化.YAMLで使うときは@を外す
const (
	// Scriptで書き換え不可
	PROP_ID    = "@ID"    // FindByIDで参照できる
	PROP_EVENT = "@Event" // ScriptでDispatch中のEvnetを参照できる
	PROP_FRAME = "@Frame" // Scriptで現在のFrameを参照できる

	// 共通
	PROP_TYPE       = "@Type"      // UIの種類
	PROP_IS_ENABLE  = "@IsEnable"  // Dispatchの対象にするかどうか
	PROP_IS_VISIBLE = "@IsVisible" // Drawの対象にするかどうか
	PROP_REMOVE     = "@Remove"    // Updateの後に削除するかどうか
	PROP_X          = "@X"         // DrawのときのX座標。親からのオフセット
	PROP_Y          = "@Y"         // DrawのときのY座標。親からのオフセット
	PROP_W          = "@W"         // Drawのときの幅.Clip用
	PROP_H          = "@H"         // Drawのときの高さ.Clip用

	// テキスト系
	PROP_TEXT = "@Text" // LabelやButtonのテキスト

	// 選択系
	PROP_SELECT_NO = "@SelectNo" // Selectの選択番号。0から始まる

	// レイアウト用
	PROP_MARGIN        = "@Margin"       // MarginTop, MarginBottom, MarginLeft, MarginRightを一括で指定するときに使う
	PROP_MARGIN_TOP    = "@MarginTop"    // MarginTop
	PROP_MARGIN_BOTTOM = "@MarginBottom" // MarginBottom
	PROP_MARGIN_LEFT   = "@MarginLeft"   // MarginLeft
	PROP_MARGIN_RIGHT  = "@MarginRight"  //
	PROP_MARGIN_X      = "@MarginX"      // MarginLeftとMarginRightを一括で指定するときに使う
	PROP_MARGIN_Y      = "@MarginY"      // MarginTopとMarginBottomを一括で指定するときに使う

	PROP_ALIGN_CENTER  = "@AlignCenter"  // AlignCenterXとAlignCenterYを一括で指定するときに使う
	PROP_ALIGN_CENTERX = "@AlignCenterX" // AlignCenterX
	PROP_ALIGN_CENTERY = "@AlignCenterY" // AlignCenterY
	PROP_ALIGN_RIGHT   = "@AlignRight"   // AlignRight
	PROP_ALIGN_BOTTOM  = "@AlignBottom"  // AlignBottom

	PROP_IS_ABS = "@IsAbs" // XやYを絶対座標で指定するかどうか
)
