package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"

	generalview "github.com/Ceinl/Go-dashboard/internal/generalView"
	"github.com/Ceinl/Go-dashboard/internal/module"
	"github.com/Ceinl/Go-dashboard/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AppConfig struct {
	LastActiveWorkspaceID string `json:"last_active_workspace_id"`
}

const (
	listState uint = iota
	projectState
	workspaceState
	dashboardState
	SettingsState
	CreateWorkspaceState
	DeleteWorkspaceState
	SwapWorkspaceState
	CreateProjectState
	ModuleSelectorState
	HelpState
	ConfirmationState
	WorkspaceModuleSelectorState
)

type model struct {
	state                       uint
	currentWorkspace            storage.Workspace
	currentProject              storage.Project
	projects                    []storage.Project
	currentModule               module.Module
	activeModules               []module.Module
	currentModuleIndex          int
	createWorkspaceView         generalview.CreateWorkspaceView
	deleteWorkspaceView         generalview.DeleteWorkspaceView
	swapWorkspaceView           generalview.SwapWorkspaceView
	createProjectView           generalview.CreateProjectView
	moduleSelectorView          generalview.ModuleSelectorView
	workspaceModuleSelectorView generalview.WorkspaceModuleSelectorView
	helpView                    generalview.HelpView
	confirmationView            generalview.ConfirmationView

	db     *sql.DB
	config AppConfig

	projectBar *generalview.ProjectBar
	statusBar  generalview.StatusBar
	body       generalview.Body

	sizeInitialized bool
	width           int
	height          int
}

func (m *model) Init() tea.Cmd {
	if m.config.LastActiveWorkspaceID != "" {
		ws, err := storage.GetWorkspace(m.db, m.config.LastActiveWorkspaceID)
		if err == nil {
			m.currentWorkspace = ws
			m.statusBar.ActiveWorkspace = ws.Name
		}
	}

	if m.currentWorkspace.ID == "" {
		workspaces, err := storage.GetAllWorkspaces(m.db)
		if err != nil {
			log.Printf("Error getting all workspaces: %v", err)
		} else if len(workspaces) > 0 {
			m.currentWorkspace = workspaces[0]
			m.statusBar.ActiveWorkspace = workspaces[0].Name
		}
	}

	m.reloadProjects()

	return tea.Batch(
		m.projectBar.Init(),
		m.statusBar.Init(),
		m.body.Init(),
		m.createWorkspaceView.Init(),
		m.deleteWorkspaceView.Init(),
		m.swapWorkspaceView.Init(),
	)
}

type YesDeleteWorkspaceMsg struct{ ID string }
type NoDeleteWorkspaceMsg struct{}
type YesDeleteProjectMsg struct{ ID string }
type NoDeleteProjectMsg struct{}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.sizeInitialized {
			cmd = m.reloadActiveModules()
			cmds = append(cmds, cmd)
			m.sizeInitialized = true
		}

		m.projectBar, cmd = m.projectBar.Update(msg)
		cmds = append(cmds, cmd)
		m.statusBar, cmd = m.statusBar.Update(msg)
		cmds = append(cmds, cmd)
		m.body, _ = m.body.Update(msg)
		m.createWorkspaceView, _ = m.createWorkspaceView.Update(msg)
		m.deleteWorkspaceView, _ = m.deleteWorkspaceView.Update(msg)
		m.swapWorkspaceView, _ = m.swapWorkspaceView.Update(msg)
		m.createProjectView, _ = m.createProjectView.Update(msg)
		if m.currentModule != nil {
			m.currentModule, cmd = m.currentModule.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		// Global quit
		if msg.String() == "alt+q" {
			return m, tea.Quit
		}

		if m.state == WorkspaceModuleSelectorState {
			m.workspaceModuleSelectorView, cmd = m.workspaceModuleSelectorView.Update(msg)
			return m, cmd
		}

		// Handle command mode exclusively
		if m.statusBar.CommandMode {
			m.statusBar, cmd = m.statusBar.Update(msg)
			return m, cmd
		}

		// Handle module switching
		if msg.String() == "shift+up" {
			if len(m.activeModules) > 0 {
				m.currentModuleIndex--
				if m.currentModuleIndex < 0 {
					m.currentModuleIndex = len(m.activeModules) - 1
				}
				m.currentModule = m.activeModules[m.currentModuleIndex]
			}
			return m, nil
		}

		if msg.String() == "shift+down" {
			if len(m.activeModules) > 0 {
				m.currentModuleIndex++
				if m.currentModuleIndex >= len(m.activeModules) {
					m.currentModuleIndex = 0
				}
				m.currentModule = m.activeModules[m.currentModuleIndex]
			}
			return m, nil
		}

		// Handle wizard states
		switch m.state {
		case CreateWorkspaceState:
			m.createWorkspaceView, cmd = m.createWorkspaceView.Update(msg)
			cmds = append(cmds, cmd)
		case DeleteWorkspaceState:
			m.deleteWorkspaceView, cmd = m.deleteWorkspaceView.Update(msg)
			cmds = append(cmds, cmd)
		case SwapWorkspaceState:
			m.swapWorkspaceView, cmd = m.swapWorkspaceView.Update(msg)
			cmds = append(cmds, cmd)
		case CreateProjectState:
			m.createProjectView, cmd = m.createProjectView.Update(msg)
			cmds = append(cmds, cmd)
		case ModuleSelectorState:
			m.moduleSelectorView, cmd = m.moduleSelectorView.Update(msg)
			cmds = append(cmds, cmd)
		case HelpState:
			m.helpView, cmd = m.helpView.Update(msg)
			cmds = append(cmds, cmd)
		case ConfirmationState:
			m.confirmationView, cmd = m.confirmationView.Update(msg)
			cmds = append(cmds, cmd)
		default:
			// If not in a wizard state, pass keys to the status bar (for entering command mode)
			// and the active module.
			m.statusBar, cmd = m.statusBar.Update(msg)
			cmds = append(cmds, cmd)

			if !m.statusBar.CommandMode && m.currentModule != nil {
				m.currentModule, cmd = m.currentModule.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

		// Always update the status bar text
		m.statusBar.ActiveWorkspace = m.currentWorkspace.Name
		m.statusBar.ActiveProject = m.currentProject.Name
	case generalview.SwitchProjectMsg:
		m.currentProject = msg.Project
		m.statusBar.ActiveProject = m.currentProject.Name
		m.reloadProjects()
		cmd = m.reloadActiveModules()
		cmds = append(cmds, cmd)
	case generalview.DoneCreateWorkspaceMsg:
		m.state = workspaceState
		m.createWorkspaceView = generalview.NewCreateWorkspaceView(m.db)
		return m, nil
	case generalview.DoneDeleteWorkspaceMsg:
		m.state = workspaceState
		m.deleteWorkspaceView = generalview.NewDeleteWorkspaceView(m.db)
		return m, nil
	case generalview.ConfirmDeleteWorkspaceMsg:
		m.state = ConfirmationState
		m.confirmationView = generalview.NewConfirmationView(
			"Are you sure you want to delete this workspace?",
			YesDeleteWorkspaceMsg{ID: msg.WorkspaceID},
			NoDeleteWorkspaceMsg{},
		)
		return m, nil
	case YesDeleteWorkspaceMsg:
		storage.DeleteWorkspace(m.db, msg.ID)
		m.state = workspaceState
		m.reloadProjects()
		return m, nil
	case NoDeleteWorkspaceMsg:
		m.state = workspaceState
		return m, nil
	case YesDeleteProjectMsg:
		storage.DeleteProject(m.db, msg.ID)
		m.state = projectState
		m.reloadProjects()
		return m, nil
	case NoDeleteProjectMsg:
		m.state = projectState
		return m, nil
	case generalview.DoneSwapWorkspaceMsg:
		m.state = workspaceState
		m.swapWorkspaceView = generalview.NewSwapWorkspaceView(m.db)
		if msg.SelectedWorkspace.ID != "" {
			m.currentWorkspace = msg.SelectedWorkspace
			m.statusBar.ActiveWorkspace = m.currentWorkspace.Name
			m.config.LastActiveWorkspaceID = m.currentWorkspace.ID
			if err := saveConfig(m.config); err != nil {
				log.Printf("Error saving config: %v", err)
			}
			m.reloadProjects()
			cmd = m.reloadActiveModules()
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case generalview.DoneCreateProjectMsg:
		m.state = projectState
		m.createProjectView = generalview.NewCreateProjectView(m.db, m.currentWorkspace.ID)
		m.reloadProjects()
		return m, nil
	case generalview.DoneModuleSelectorMsg:
		m.state = projectState
		m.currentProject = msg.Project
		storage.UpdateProject(m.db, m.currentProject)
		m.reloadProjects()
		return m, nil
	case generalview.DoneWorkspaceModuleSelectorMsg:
		m.state = projectState
		m.currentWorkspace = msg.Workspace
		storage.UpdateWorkspace(m.db, m.currentWorkspace)
		cmd = m.reloadActiveModules()
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case generalview.DoneHelpMsg:
		m.state = projectState
		return m, nil
	case generalview.NewWorkspaceCommandMsg:
		m.state = CreateWorkspaceState
		m.createWorkspaceView = generalview.NewCreateWorkspaceView(m.db)
		cmds = append(cmds, m.createWorkspaceView.Init())
	case generalview.DeleteWorkspaceCommandMsg:
		m.state = DeleteWorkspaceState
		m.deleteWorkspaceView = generalview.NewDeleteWorkspaceView(m.db)
		cmds = append(cmds, m.deleteWorkspaceView.Init())
	case generalview.DeleteProjectCommandMsg:
		m.state = ConfirmationState
		m.confirmationView = generalview.NewConfirmationView(
			"Are you sure you want to delete this project?",
			YesDeleteProjectMsg{ID: m.currentProject.ID},
			NoDeleteProjectMsg{},
		)
		return m, nil
	case generalview.SwapWorkspaceCommandMsg:
		m.state = SwapWorkspaceState
		m.swapWorkspaceView = generalview.NewSwapWorkspaceView(m.db)
		cmds = append(cmds, m.swapWorkspaceView.Init())
	case generalview.NewProjectCommandMsg:
		m.state = CreateProjectState
		m.createProjectView = generalview.NewCreateProjectView(m.db, m.currentWorkspace.ID)
		cmds = append(cmds, m.createProjectView.Init())
	case generalview.ModuleSelectorCommandMsg:
		m.state = ModuleSelectorState
		m.moduleSelectorView = generalview.NewModuleSelectorView(m.currentProject)
		cmds = append(cmds, m.moduleSelectorView.Init())
	case generalview.WorkspaceModuleSelectorCommandMsg:
		m.state = WorkspaceModuleSelectorState
		m.workspaceModuleSelectorView = generalview.NewWorkspaceModuleSelectorView(m.currentWorkspace)
		cmds = append(cmds, m.workspaceModuleSelectorView.Init())
	case generalview.HelpCommandMsg:
		m.state = HelpState
		m.helpView = generalview.NewHelpView()
		cmds = append(cmds, m.helpView.Init())
	}

	if m.state != CreateWorkspaceState && m.state != DeleteWorkspaceState && m.state != SwapWorkspaceState && m.state != CreateProjectState && m.state != ModuleSelectorState && m.state != HelpState && m.state != ConfirmationState && m.state != WorkspaceModuleSelectorState {
		var projectBarCmd tea.Cmd
		m.projectBar, projectBarCmd = m.projectBar.Update(msg)
		cmds = append(cmds, projectBarCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if !m.sizeInitialized {
		return "Loading..."
	}

	// Handle wizard states first
	if m.state == CreateWorkspaceState {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.createWorkspaceView.View())
	} else if m.state == DeleteWorkspaceState {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.deleteWorkspaceView.View())
	} else if m.state == SwapWorkspaceState {
		return m.swapWorkspaceView.View()
	} else if m.state == CreateProjectState {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.createProjectView.View())
	} else if m.state == ModuleSelectorState {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.moduleSelectorView.View())
	} else if m.state == WorkspaceModuleSelectorState {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.workspaceModuleSelectorView.View())
	} else if m.state == HelpState {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.helpView.View())
	} else if m.state == ConfirmationState {
		return m.confirmationView.View()
	}

	// Regular view layout
	projectBarView := m.projectBar.View()
	statusBarView := m.statusBar.View()

	availableHeight := m.height - lipgloss.Height(projectBarView) - lipgloss.Height(statusBarView)

	var middleView string
	if m.currentModule != nil {
		middleView = m.currentModule.View()
	} else {
		middleView = "" // No active module
	}

	// Create a container for the middle view that will handle sizing and centering.
	middleContainer := lipgloss.NewStyle().
		Height(availableHeight).
		Width(m.width).
		Align(lipgloss.Center).
		Render(middleView)

	return lipgloss.JoinVertical(lipgloss.Left,
		projectBarView,
		middleContainer,
		statusBarView,
	)
}

func (m *model) reloadProjects() {
	projects, err := storage.GetAllProjectsForWorkspace(m.db, m.currentWorkspace.ID)
	if err != nil {
		log.Printf("Error getting all projects: %v", err)
		m.projects = []storage.Project{}
	} else {
		m.projects = projects
	}

	if len(m.projects) > 0 {
		// If the current project is no longer in the list, reset it
		found := false
		for i, p := range m.projects {
			if p.ID == m.currentProject.ID {
				m.projectBar.SelectedIndex = i
				found = true
				break
			}
		}
		if !found {
			m.currentProject = m.projects[0]
			m.projectBar.SelectedIndex = 0
		}
	} else {
		m.currentProject = storage.Project{}
		m.projectBar.SelectedIndex = -1
	}

	m.projectBar.Projects = m.projects
	m.statusBar.ActiveProject = m.currentProject.Name
}

func (m *model) reloadActiveModules() tea.Cmd {
	var initCmds []tea.Cmd
	m.activeModules = []module.Module{}
	if m.currentWorkspace.ID != "" && m.currentProject.ID != "" {
		moduleNames := strings.Split(m.currentWorkspace.ActiveModules, ",")
		for _, name := range moduleNames {
			if name == "" {
				continue
			}
			var newModule module.Module
			switch name {
			case "linksaver":
				newModule = module.NewLinkSaver(m.db, m.currentProject.ID)
			case "placeholder":
				newModule = module.NewPlaceholder()
			case "kanban":
				newModule = module.NewKanban(m.db, m.currentProject.ID)
			case "twitter":
				newModule = module.NewTwitter(m.db, m.currentProject.ID)
			}
			if newModule != nil {
				if m.width > 0 && m.height > 0 {
					var cmd tea.Cmd
					newModule, cmd = newModule.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
					initCmds = append(initCmds, cmd)
				}
				m.activeModules = append(m.activeModules, newModule)
				initCmds = append(initCmds, newModule.Init())
			}
		}
	}

	if len(m.activeModules) > 0 {
		if m.currentModuleIndex >= len(m.activeModules) {
			m.currentModuleIndex = 0
		}
		m.currentModule = m.activeModules[m.currentModuleIndex]
	} else {
		m.currentModule = nil
	}
	return tea.Batch(initCmds...)
}

func loadConfig() (AppConfig, error) {
	var config AppConfig
	file, err := os.Open("settings.json")
	if err != nil {
		if os.IsNotExist(err) {
			return AppConfig{}, nil
		}
		return AppConfig{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func saveConfig(config AppConfig) error {
	file, err := os.Create("settings.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
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

	config, err := loadConfig()
	if err != nil {
		log.Printf("Error loading config: %v. Using defaults.", err)
	}

	projectBar := generalview.NewProjectBar()
	initialModel := model{
		db:                  db,
		config:              config,
		createWorkspaceView: generalview.NewCreateWorkspaceView(db),
		deleteWorkspaceView: generalview.NewDeleteWorkspaceView(db),
		swapWorkspaceView:   generalview.NewSwapWorkspaceView(db),
		createProjectView:   generalview.NewCreateProjectView(db, ""),
		moduleSelectorView:  generalview.NewModuleSelectorView(storage.Project{}),
		helpView:            generalview.NewHelpView(),
		projectBar:          &projectBar,
	}

	p := tea.NewProgram(&initialModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
