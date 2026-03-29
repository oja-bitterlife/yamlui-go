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
	root   *yamlui.UIBase
	width  int
	height int
	frame  int

	canvas [][]Cell
	label1 *BTLabel
}

func initialModel() model {
	model := model{
		root:   yamlui.NewUIBase("Root"),
		width:  64,
		height: 24,
	}

	area := yamlui.NewUIArea()
	model.root.AddChild(area.Base)

	model.label1 = NewBTLabel("Hello, YAMLUI!")
	area.Base.AddChild(model.label1.Base.Base)

	area.Base.X = 10
	area.Base.Y = 3

	// Lispスクリプト: 実行するたびにX座標を増やし、テキストを書き換える
	scriptSrc := `
(set @X
	(switch (> @X 30)
		(set @X 0)
		(+ @X 1)))
	   `
	if err := model.label1.Base.Base.SetScript(scriptSrc); err != nil {
		panic(fmt.Sprintf("Failed to set script: %v", err))
	}

	model.canvas = make([][]Cell, model.width)
	for i := range model.canvas {
		// 内側のスライス（列）を作る
		model.canvas[i] = make([]Cell, model.width)
		// 初期状態としてスペースなどで埋める
		for j := range model.canvas[i] {
			model.canvas[i][j] = Cell{Rune: 'x', Color: "white"}
		}
	}

	return model
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

		// キーが押されたらLispを実行して構造体を更新
		err := m.root.UpdateTree(m.frame)
		if err != nil {
			fmt.Printf("Error executing script: %v\n", err)
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {

	m.label1.Canvas = m.canvas

	m.root.DrawTree(0, 0) // 描画用の構造体を更新

	// ラベルに色と位置（マージン）を適用
	// style := lipgloss.NewStyle().
	// 	Foreground(m.getColor(m.root.Color)).
	// 	MarginLeft(m.root.X).
	// 	MarginTop(m.root.Y)

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

	// return style.Render(m.root.Text)
	return b.String()
}

func (m model) getColor(c string) lipgloss.Color {
	if c == "system" {
		return lipgloss.Color("86") // 例えば少し青みがかったシアンなど
	}
	return lipgloss.Color(c)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
