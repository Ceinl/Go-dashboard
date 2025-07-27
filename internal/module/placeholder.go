package module

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Placeholder struct {
	width  int
	height int
}

func NewPlaceholder() Module {
	return &Placeholder{}
}

func (m *Placeholder) Init() tea.Cmd {
	return nil
}

func (m *Placeholder) Update(msg tea.Msg) (Module, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *Placeholder) View() string {
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render("This is a placeholder module.")
}
