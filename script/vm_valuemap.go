package script

// **********************************************************************
// Mapから値を拾うのが思ったより面倒なので、ValueMapを作成してメッソドを生やす
type ValueMap map[string]Value

// ==================================================
// 型ごとに取得する関数を用意。存在しない場合や型が違う場合はデフォルト値を返す
func (data ValueMap) GetStr(key string, def string) string {
	value, ok := data[key]
	if !ok || value.Type != TypeString {
		return def
	}
	return value.Str
}

func (data ValueMap) GetNum(key string, def int) int {
	value, ok := data[key]
	if !ok || value.Type != TypeNumber {
		return def
	}
	return value.Num
}

func (data ValueMap) GetBool(key string, def bool) bool {
	value, ok := data[key]
	if !ok || value.Type != TypeBool {
		return def
	}
	return value.Bool
}
