package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
	"github.com/oja-bitterlife/yamlui-go/yamlui_json"
)

type Cell struct {
	Rune  rune
	Color string
}

type model struct {
	lib    *yamlui.YAMLUI
	width  int
	height int
	frame  int

	canvas   [][]Cell
	startSel *BTSpeed
	speedSel *BTSpeed
}

func initialModel() model {
	m := model{
		lib:    yamlui.NewYAMLUI(),
		width:  68,
		height: 26,
	}
	m.canvas = make([][]Cell, m.height)
	for i := range m.canvas {
		// 内側のスライス（列）を作る
		m.canvas[i] = make([]Cell, m.width)
		// 初期状態としてスペースなどで埋める
		for j := range m.canvas[i] {
			m.canvas[i][j] = Cell{Rune: ' ', Color: "white"}
		}
	}

	// UI構築の登録
	m.lib.RegisterRemap("window", func(type_ string, parent *yamlui.UIBase, data map[string]script.Value) (*yamlui.UIBase, error) {
		return NewBTWindow(&m, parent, data).Base.Base, nil
	})
	m.lib.RegisterRemap("title", func(type_ string, parent *yamlui.UIBase, data map[string]script.Value) (*yamlui.UIBase, error) {
		return NewBTTitle(&m, parent, data).Base.Base, nil
	})
	m.lib.RegisterRemap("area", func(type_ string, parent *yamlui.UIBase, data map[string]script.Value) (*yamlui.UIBase, error) {
		return yamlui.NewUIArea(parent, data).Base, nil
	})
	m.lib.RegisterRemap("start", func(type_ string, parent *yamlui.UIBase, data map[string]script.Value) (*yamlui.UIBase, error) {
		return NewBTStart(&m, parent, data).Base.Base, nil
	})
	m.lib.RegisterRemap("label", func(type_ string, parent *yamlui.UIBase, data map[string]script.Value) (*yamlui.UIBase, error) {
		return NewBTLabel(&m, parent, data).Base.Base, nil
	})
	m.lib.RegisterRemap("speed", func(type_ string, parent *yamlui.UIBase, data map[string]script.Value) (*yamlui.UIBase, error) {
		return NewBTSpeed(&m, parent, data).Base.Base, nil
	})

	// JSON を読み込んで UI を構築する
	fileData, err := os.ReadFile("bin/ui.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to read ui.json: %v", err))
	}
	data, err := yamlui_json.AnyJSONToValueJSON(fileData)

	if err := m.lib.Load(data); err != nil {
		panic(fmt.Sprintf("Failed to load UI from JSON: %v", err))
	}

	// Lispスクリプト: 実行するたびにX座標を増やし、テキストを書き換える
	// 	scriptSrc := `
	// (set @X
	// 	(switch (< @X 3)
	// 		(+ @X 1)
	// 		(set @X 0)))
	// 	   `
	// 	if err := label.Base.Base.SetScript(scriptSrc); err != nil {
	// 		panic(fmt.Sprintf("Failed to set script: %v", err))
	// 	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.frame++ // フレームを進める

	m.lib.ClearEvents() // イベントをクリア

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 'q' か 'ctrl+c' で終了
		if msg.String() == "q" || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// 縦方向の移動 (NextGridY)
		if msg.Type == tea.KeyUp {
			m.lib.AddEvent("key:up")
		}
		if msg.Type == tea.KeyDown {
			m.lib.AddEvent("key:down")
		}

		// // 横方向の移動 (NextGridX)
		if msg.Type == tea.KeyLeft {
			m.lib.AddEvent("key:left")
		}
		if msg.Type == tea.KeyRight {
			m.lib.AddEvent("key:right")
		}
		// if msg.Type == tea.KeyLeft {
		// 	m.speedSel.Base.NextGridX(-1, selectToggle)
		// }
		// if msg.Type == tea.KeyRight {
		// 	m.speedSel.Base.NextGridX(1, selectToggle)
		// }
	}

	// Update
	errorList := m.lib.Update(m.frame)
	for _, err := range errorList {
		fmt.Printf("Error during update: %v\n", err)
	}

	return m, nil
}

func (m model) View() string {
	for y := range m.canvas {
		for x := range m.canvas[y] {
			m.canvas[y][x] = Cell{Rune: ' ', Color: "white"}
		}
	}

	m.lib.Draw(yamlui.NewAreaI(0, 0, m.width, m.height))

	var b strings.Builder
	for y := 0; y < len(m.canvas); y++ {
		for x := 0; x < len(m.canvas[y]); x++ {
			cell := m.canvas[y][x]
			// ここで lipgloss などを使って色をつけても良いですが、
			// まずはシンプルに文字だけ出すなら：
			b.WriteRune(cell.Rune)
		}
		b.WriteByte('\n')
	}

	return b.String()
}

func (m model) getColor(c string) lipgloss.Color {
	if c == "system" {
		return lipgloss.Color("86") // 例えば少し青みがかったシアンなど
	}
	return lipgloss.Color(c)
}

func main() {
	// p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
