package script

import (
	"testing"
)

func TestValue_UnmarshalJSON(t *testing.T) {
	// 1. テストデータの準備 (数値、文字列、プロパティ、リストの混合)
	orig := Value{
		Type: TypeList,
		List: []Value{
			{Type: TypeNumber, Num: 1.2345},
			{Type: TypeString, Str: "hello"},
			{Type: TypeProperty, Str: ":color"},
			{Type: TypeBool, Bool: true},
		},
	}

	// 2. Marshal (Value -> JSON)
	data, err := orig.MarshalJSON()
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// 3. Unmarshal (JSON -> 新しい Value)
	var restored Value
	err = restored.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// 4. 検証 (地道に中身をチェック)
	if restored.Type != orig.Type {
		t.Errorf("Type mismatch: expected %v, got %v", orig.Type, restored.Type)
	}

	if len(restored.List) != len(orig.List) {
		t.Fatalf("List length mismatch: expected %d, got %d", len(orig.List), len(restored.List))
	}

	// 個別の要素をチェック
	if restored.List[0].Num != orig.List[0].Num {
		t.Errorf("List[0] Num mismatch: expected %f, got %f", orig.List[0].Num, restored.List[0].Num)
	}
	if restored.List[1].Str != orig.List[1].Str {
		t.Errorf("List[1] Str mismatch: expected %q, got %q", orig.List[1].Str, restored.List[1].Str)
	}
	if restored.List[2].Type != TypeProperty || restored.List[2].Str != ":color" {
		t.Errorf("List[2] Property mismatch")
	}
}

func TestValue_UnmarshalJSON_Spaces(t *testing.T) {
	// クレンジング（TrimSpace）が効いているか確認するテスト
	input := `  {  "Num"  :  99.9  }  `
	var v Value
	if err := v.UnmarshalJSON([]byte(input)); err != nil {
		t.Fatalf("Failed to parse spaced JSON: %v", err)
	}
	if v.Type != TypeNumber || v.Num != 99.9 {
		t.Errorf("Spaced JSON parse error: got %+v", v)
	}
}
