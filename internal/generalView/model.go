package generalview

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Workspace struct {
	ID    string
	Name  string
	Color string
}

type Model struct {
	projectBar ProjectBar
	statusBar  StatusBar
}

func NewModel() Model {
	return Model{
		projectBar: ProjectBar{},
		statusBar:  StatusBar{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.projectBar.Init(),
		m.statusBar.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var pbCmd tea.Cmd
	m.projectBar, pbCmd = m.projectBar.Update(msg)
	cmds = append(cmds, pbCmd)

	var sbCmd tea.Cmd
	m.statusBar, sbCmd = m.statusBar.Update(msg)
	cmds = append(cmds, sbCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	height := 24
	bodyHeight := height - 2

	return strings.Join([]string{
		centerText(m.projectBar.View(), 80),
		strings.Repeat("\n", bodyHeight),
		centerText(m.statusBar.View(), 80),
	}, "\n")
}

func centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text
}
