package main

import (
	"database/sql"
	"log"

	generalview "github.com/Ceinl/Go-dashboard/internal/generalView"
	"github.com/Ceinl/Go-dashboard/internal/module"
	"github.com/Ceinl/Go-dashboard/internal/storage"
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
	DeleteWorkspaceState
	SwapWorkspaceState
)

type model struct {
	state               uint
	modules             []module.Module
	workspace           workspace.Workspace
	currentWorkspace    storage.Workspace
	createWorkspaceView generalview.CreateWorkspaceView
	deleteWorkspaceView generalview.DeleteWorkspaceView
	swapWorkspaceView   generalview.SwapWorkspaceView

	db *sql.DB

	projectBar generalview.ProjectBar
	statusBar  generalview.StatusBar
	body       generalview.Body
	width      int
	height     int
}

func (m model) Init() tea.Cmd {
	// Initialize currentWorkspace with a default or first available workspace
	workspaces, err := storage.GetAllWorkspaces(m.db)
	if err != nil {
		log.Printf("Error getting all workspaces: %v", err)
	} else if len(workspaces) > 0 {
		m.currentWorkspace = workspaces[0]
	}

	
	return tea.Batch(
		m.projectBar.Init(),
		m.statusBar.Init(),
		m.body.Init(),
		m.createWorkspaceView.Init(),
		m.deleteWorkspaceView.Init(),
		m.swapWorkspaceView.Init(),
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
		m.deleteWorkspaceView, _ = m.deleteWorkspaceView.Update(msg)
		m.swapWorkspaceView, _ = m.swapWorkspaceView.Update(msg)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "alt+q" {
			return m, tea.Quit
		}

		if m.state == CreateWorkspaceState {
			var cmd tea.Cmd
			m.createWorkspaceView, cmd = m.createWorkspaceView.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		} else if m.state == DeleteWorkspaceState {
			var cmd tea.Cmd
			m.deleteWorkspaceView, cmd = m.deleteWorkspaceView.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		} else if m.state == SwapWorkspaceState {
			var cmd tea.Cmd
			m.swapWorkspaceView, cmd = m.swapWorkspaceView.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		} else {
			var cmd tea.Cmd
			m.statusBar, cmd = m.statusBar.Update(msg)
			m.statusBar.ActiveWorkspace = m.currentWorkspace.Name
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			if m.statusBar.CommandMode {
				return m, tea.Batch(cmds...)
			}
		}
	case generalview.DoneCreateWorkspaceMsg:
		m.state = workspaceState
		m.createWorkspaceView = generalview.NewCreateWorkspaceView(m.db)
		return m, nil
	case generalview.DoneDeleteWorkspaceMsg:
		m.state = workspaceState
		m.deleteWorkspaceView = generalview.NewDeleteWorkspaceView(m.db)
		return m, nil
	case generalview.DoneSwapWorkspaceMsg:
		m.state = workspaceState
		m.swapWorkspaceView = generalview.NewSwapWorkspaceView(m.db)
		if msg.SelectedWorkspace.ID != "" {
			m.currentWorkspace = msg.SelectedWorkspace
		}
		return m, nil
	case generalview.NewWorkspaceCommandMsg:
		m.state = CreateWorkspaceState
		m.createWorkspaceView = generalview.NewCreateWorkspaceView(m.db)
		cmds = append(cmds, m.createWorkspaceView.Init())
	case generalview.DeleteWorkspaceCommandMsg:
		m.state = DeleteWorkspaceState
		m.deleteWorkspaceView = generalview.NewDeleteWorkspaceView(m.db)
		cmds = append(cmds, m.deleteWorkspaceView.Init())
	case generalview.SwapWorkspaceCommandMsg:
		m.state = SwapWorkspaceState
		m.swapWorkspaceView = generalview.NewSwapWorkspaceView(m.db)
		cmds = append(cmds, m.swapWorkspaceView.Init())
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
	} else if m.state == DeleteWorkspaceState {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.deleteWorkspaceView.View())
	} else if m.state == SwapWorkspaceState {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.swapWorkspaceView.View())
	}

	return mainView
}

func main() {
	f, err := tea.LogToFile("debug.Log", "debug")
	if err != nil {
		log.Fatalf("err: %w", err)
	}
	defer f.Close()

	db, err := storage.InitDB("file:test.db?_foreign_keys=on")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	initialModel := model{
		db:                  db,
		createWorkspaceView: generalview.NewCreateWorkspaceView(db),
		deleteWorkspaceView: generalview.NewDeleteWorkspaceView(db),
		swapWorkspaceView:   generalview.NewSwapWorkspaceView(db),
	}

	p := tea.NewProgram(&initialModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
