package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	// m.lib.UIBuild("window", NewBTWindow(lib.model))
	// yamlui.UIBuilder(m.lib, "window", NewBTWindow, func(ui *BTWindow) { ui.model = &m })
	// yamlui.UIBuilder(m.lib, "title", NewBTTitle, func(ui *BTTitle) { ui.model = &m })
	// yamlui.UIBuilder(m.lib, "area", yamlui.NewUIArea, nil)
	// yamlui.UIBuilder(m.lib, "start", NewBTStart, func(ui *BTStart) { ui.model = &m })
	// yamlui.UIBuilder(m.lib, "label", NewBTLabel, func(ui *BTLabel) { ui.model = &m })
	// yamlui.UIBuilder(m.lib, "speed", NewBTSpeed, func(ui *BTSpeed) { ui.model = &m })

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
			m.canvas[y][x] = Cell{Rune: ' ', Color: "navy"}
		}
	}

	m.lib.Draw(yamlui.NewAreaI(0, 0, m.width, m.height))

	// 色の設定
	bgStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#000080")). // 紺色背景
		Foreground(lipgloss.Color("#FFFFFF"))  // 白文字

	var b strings.Builder
	for y := range m.canvas {
		var row strings.Builder
		for x := range m.canvas[y] {
			row.WriteRune(m.canvas[y][x].Rune)
		}
		// 1行まとめて紺色背景にする
		line := bgStyle.Render(row.String())
		b.WriteString(line + "\n")
	}

	return b.String()
}

func main() {
	// p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
