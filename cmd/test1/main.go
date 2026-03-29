package main

import (
	"fmt"
	"os"

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
}

func initialModel() model {
	area := yamlui.NewUIArea()
	label := yamlui.NewUILabel("Hello, YAMLUI!")
	area.Base.AddChild(label.Base)

	area.Base.X = 10
	area.Base.Y = 3

	// Lispスクリプト: 実行するたびにX座標を増やし、テキストを書き換える
	scriptSrc := `
(set @X
	(switch (> @X 30)
		(set @X 0)
		(+ @X 1)))
	   `
	if err := label.Base.SetScript(scriptSrc); err != nil {
		panic(fmt.Sprintf("Failed to set script: %v", err))
	}

	root := yamlui.NewUIBase("Root")
	root.AddChild(area.Base)
	return model{root: root}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 'q' か 'ctrl+c' で終了
		if msg.String() == "q" || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// キーが押されたらLispを実行して構造体を更新
		err := m.root.Update()
		if err != nil {
			fmt.Printf("Error executing script: %v\n", err)
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	// ラベルに色と位置（マージン）を適用
	style := lipgloss.NewStyle().
		Foreground(m.getColor(m.root.Color)).
		MarginLeft(m.root.X).
		MarginTop(m.root.Y)

	return style.Render(m.root.Text)
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
