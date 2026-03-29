package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oja-bitterlife/yamlui-go/yamlui"
)

type model struct {
	label *yamlui.UILabelBase
}

func initialModel() model {
	label := yamlui.NewUILabelBase()
	label.Base.Text = "Hello, YAMLUI!"

	// Lispスクリプト: 実行するたびにX座標を増やし、テキストを書き換える
	scriptSrc := `
        (set @X (+ @X 1))
        (if (> @X 30) (set @X 0))
    `
	label.Base.SetScript(scriptSrc)

	return model{label: label}
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
		m.label.Base.Update()
	}

	return m, nil
}

func (m model) View() string {
	// 座標に合わせて描画
	var out strings.Builder
	for i := 0; i < m.label.Base.Y; i++ {
		out.WriteString("\n")
	}

	var indent strings.Builder
	for i := 0; i < m.label.Base.X; i++ {
		indent.WriteString(" ")
	}

	return fmt.Sprintf("%s%s%s\n\n(x:%d, count:%d) Press 'q' to quit, other keys to move",
		out.String(), indent.String(), m.label.Base.Text, m.label.Base.X, m.label.Base.Frame)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
