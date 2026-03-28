package main

import (
	"encoding/json"

	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type UI struct {
	name string
}

func main() {
	yamlui.NewYamlUI(
		func(data *UI) ([]byte, error) {
			return json.Marshal(data)
		},
		func(b []byte) (*UI, error) {
			var data UI
			err := json.Unmarshal(b, &data)
			return &data, err
		},
	)
}
