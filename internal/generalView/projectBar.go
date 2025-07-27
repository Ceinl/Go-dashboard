package generalview

import (
	"github.com/Ceinl/Go-dashboard/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProjectBar struct {
	Projects      []storage.Project
	SelectedIndex int
	width         int
}

func NewProjectBar() ProjectBar {
	return ProjectBar{}
}

func (m *ProjectBar) Init() tea.Cmd {
	return nil
}

type SwitchProjectMsg struct {
	Project storage.Project
}

func (m *ProjectBar) Update(msg tea.Msg) (*ProjectBar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "shift+left", "shift+h":
			if len(m.Projects) > 0 {
				m.SelectedIndex--
				if m.SelectedIndex < 0 {
					m.SelectedIndex = len(m.Projects) - 1
				}
				return m, func() tea.Msg {
					return SwitchProjectMsg{Project: m.Projects[m.SelectedIndex]}
				}
			}
		case "shift+right", "shift+l":
			if len(m.Projects) > 0 {
				m.SelectedIndex++
				if m.SelectedIndex >= len(m.Projects) {
					m.SelectedIndex = 0
				}
				return m, func() tea.Msg {
					return SwitchProjectMsg{Project: m.Projects[m.SelectedIndex]}
				}
			}
		}
	}
	return m, nil
}

func (m *ProjectBar) View() string {
	var tabs []string
	for i, p := range m.Projects {
		style := lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("255"))
		if i == m.SelectedIndex {
			style = style.Bold(true).
				Background(lipgloss.Color("57")).
				Foreground(lipgloss.Color("255"))
		}
		tabs = append(tabs, style.Render(p.Name))
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
}