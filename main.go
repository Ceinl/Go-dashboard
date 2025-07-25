package main

import (
	"log"

	generalview "github.com/Ceinl/Go-dashboard/internal/generalView"
	"github.com/Ceinl/Go-dashboard/internal/module"
	"github.com/Ceinl/Go-dashboard/internal/workspace"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listState uint = iota
	projectState
	workspaceState
	dashboardState
	SettingsState
	CreateProjectState
	CreateWorkspaceState
)

type model struct {
	state               uint
	modules             []module.Module
	workspace           workspace.Workspace
	createWorkspaceView generalview.CreateWorkspaceView

	projectBar generalview.ProjectBar
	statusBar  generalview.StatusBar
	body       generalview.Body
	width      int
	height     int
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.projectBar.Init(),
		m.statusBar.Init(),
		m.body.Init(),
		m.createWorkspaceView.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.projectBar, _ = m.projectBar.Update(msg)
		m.statusBar, _ = m.statusBar.Update(msg)
		m.body, _ = m.body.Update(msg)
		m.createWorkspaceView, _ = m.createWorkspaceView.Update(msg)
	case tea.KeyMsg:
		if msg.String() == "alt+q" {
			return m, tea.Quit
		}
		if m.state == CreateWorkspaceState {
			var cmd tea.Cmd
			m.createWorkspaceView, cmd = m.createWorkspaceView.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.statusBar, cmd = m.statusBar.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	case generalview.DoneCreateWorkspaceMsg:
		m.state = workspaceState // Or whatever state you want to return to
	case generalview.NewWorkspaceCommandMsg:
		m.state = CreateWorkspaceState
		m.createWorkspaceView = generalview.CreateWorkspaceView{}
		cmds = append(cmds, m.createWorkspaceView.Init())
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.height == 0 || m.width == 0 {
		return "Loading..."
	}

	mainView := m.projectBar.View() +
		m.body.View() +
		"\n" + m.statusBar.View()

	if m.state == CreateWorkspaceState {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.createWorkspaceView.View())
	}

	return mainView
}

func main() {
	f, err := tea.LogToFile("debug.Log", "debug")
	if err != nil {
		log.Fatalf("err: %w", err)
	}
	defer f.Close()
	p := tea.NewProgram(model{}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
