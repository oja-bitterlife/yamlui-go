package script

import (
	"encoding/json"
	"testing"
)

func TestValue_MarshalJSON(t *testing.T) {
	// 数値のテスト
	v := NewNumber(123.45678)
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	// 期待値（小数点4桁固定に設定したなら "123.4568" など）
	expected := "{\"Num\":123.4568}"
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}

	// リストのテスト
	list := NewList([]Value{NewNumber(1), NewString("hello")})
	data, err = json.Marshal(list)
	// 期待値: [1.0000,"hello"]
	t.Logf("JSON output: %s", string(data))
}
