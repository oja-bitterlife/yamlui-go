package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

const (
	TEA_UPDATE_EVENT = "tea:update"
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
	m.lib.UIBuild("close_win", NewBTCloseWin(&m))

	// JSON を読み込んで UI を構築する
	fileData, err := os.ReadFile("bin/ui.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to read ui.json: %v", err))
	}

	if err = m.lib.Load(fileData); err != nil {
		panic(fmt.Sprintf("Failed to load UI from JSON: %v", err))
	}

	TraceMemory() // 初期化時のメモリ使用状況をログに出力

	return m
}

func (m model) Init() tea.Cmd {
	return nil
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

		TraceMemory() // メモリ使用状況をログに出力
	}

	// Update
	errorList := m.lib.Update(m.frame)
	for _, err := range errorList {
		script.LogFatal("Error during update: %v", err)
	}
	if len(m.lib.GetEvents()) > 0 {
		duration := time.Duration(33 * time.Millisecond)
		return m, tea.Tick(duration, func(t time.Time) tea.Msg {
			return struct{}{} // 空のメッセージを返す
		})
	}
	return m, nil
}

func (m model) View() string {
	for y := range m.canvas {
		for x := range m.canvas[y] {
			m.canvas[y][x] = Cell{Rune: ' ', Color: "navy"}
		}
	}

	m.lib.Draw(m.width, m.height)

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
	file, err := os.Create("bin/yamlui.log")
	if err != nil {
		panic(err)
	}
	defer file.Close() // 最後に確実に閉じる
	script.SetLogWriter(file)

	// p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
