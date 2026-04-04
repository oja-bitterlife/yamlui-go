package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
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
	m.lib.UIBuild("window", NewBTWindow(&m))
	m.lib.UIBuild("title", NewBTTitle(&m))
	m.lib.UIBuild("area", yamlui.NewUILayout())
	m.lib.UIBuild("start", NewBTStart(&m))
	m.lib.UIBuild("label", NewBTLabel(&m))
	m.lib.UIBuild("speed", NewBTSpeed(&m))

	// JSON を読み込んで UI を構築する
	fileData, err := os.ReadFile("bin/ui.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to read ui.json: %v", err))
	}

	if err := m.lib.Load(fileData); err != nil {
		panic(fmt.Sprintf("Failed to load UI from JSON: %v", err))
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func nextTick() tea.Cmd {
	return func() tea.Msg {
		return struct{}{} // 空のメッセージを返す
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.frame++ // フレームを進める

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

		// Enterキーで選択
		if msg.Type == tea.KeyEnter {
			m.lib.AddEvent("key:enter")
		}
	}

	// Update
	errorList := m.lib.Update(m.frame)
	for _, err := range errorList {
		panic(fmt.Sprintf("Error during update: %v\n", err))
	}
	if len(m.lib.GetEvents()) > 0 {
		fmt.Printf("Events after update: %v\n", m.lib.GetEvents())
		return m, nextTick()
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
