package script

type ValueMap map[string]Value

func NewValueMap() ValueMap {
	return make(map[string]Value)
}

func (data ValueMap) GetStr(key string, def string) string {
	value, ok := data[key]
	if !ok || value.Type != TypeString {
		return def
	}
	return value.Str
}

func (data ValueMap) GetNum(key string, def float64) float64 {
	value, ok := data[key]
	if !ok || value.Type != TypeNumber {
		return def
	}
	return value.Num
}

func (data ValueMap) GetInt(key string, def int) int {
	value, ok := data[key]
	if !ok || value.Type != TypeNumber {
		return def
	}
	return int(value.Num)
}

func (data ValueMap) GetBool(key string, def bool) bool {
	value, ok := data[key]
	if !ok || value.Type != TypeBool {
		return def
	}
	return value.Bool
}
