package generalview

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ─────────────────────────────────────────────

type Body struct {
	Width  int
	Height int
}

func (b Body) Init() tea.Cmd {
	return nil
}

func (b Body) Update(msg tea.Msg) (Body, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.Width = msg.Width
		b.Height = msg.Height
	}
	return b, nil
}

func (b Body) View() string {
	bodyHeight := b.Height - 2
	if bodyHeight < 0 {
		bodyHeight = 0
	}
	return strings.Repeat("\n", bodyHeight)
}
