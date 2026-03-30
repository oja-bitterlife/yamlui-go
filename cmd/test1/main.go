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

	canvas   [][]Cell
	startSel *BTSelect
	speedSel *BTSelect
}

func initialModel() model {
	m := model{
		root:   yamlui.NewUIBase("Root"),
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

	win := NewBTWindow(&m)
	win.Base.Base.SetRect(0, 0, m.width, m.height)
	m.root.AddChild(win.Base.Base)

	margin := yamlui.NewUIArea()
	margin.MarginX = 2
	margin.MarginY = 1
	win.Base.Base.AddChild(margin.Base)

	title := NewBTTitle(&m)
	margin.Base.AddChild(title.Base.Base)

	startSel := NewBTSelect(&m, 1)
	startSel.Base.AddItem(yamlui.NewUISelectItem("START"))
	startSel.Base.AddItem(yamlui.NewUISelectItem("CONTINUE"))
	startSel.Base.Base.SetArea(margin.Base.Area())
	margin.Base.AddChild(startSel.Base.Base)
	m.startSel = startSel

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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 'q' か 'ctrl+c' で終了
		if msg.String() == "q" || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// キーが押されたらLispを実行して構造体を更新
		err := m.root.Update(m.frame)
		if err != nil {
			fmt.Printf("Error executing script: %v\n", err)
			return m, nil
		}

		// 縦方向の移動 (NextGridY)
		selectToggle := true
		if msg.Type == tea.KeyUp {
			m.startSel.Base.NextGridY(-1, selectToggle)
		}
		if msg.Type == tea.KeyDown {
			m.startSel.Base.NextGridY(1, selectToggle)
		}

		// 横方向の移動 (NextGridX)
		if msg.Type == tea.KeyLeft {
			m.startSel.Base.NextGridX(-1, selectToggle)
		}
		if msg.Type == tea.KeyRight {
			m.startSel.Base.NextGridX(1, selectToggle)
		}
	}

	return m, nil
}

func (m model) View() string {
	for y := range m.canvas {
		for x := range m.canvas[y] {
			m.canvas[y][x] = Cell{Rune: ' ', Color: "white"}
		}
	}

	m.root.Draw(yamlui.NewArea(0, 0, m.width, m.height))

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
