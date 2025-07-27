package generalview

import (
	"strings"

	"github.com/Ceinl/Go-dashboard/internal/module"
	"github.com/Ceinl/Go-dashboard/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WorkspaceModuleSelectorView struct {
	workspace      storage.Workspace
	availableModules []string
	cursor         int
	selected       map[string]struct{}
}

func NewWorkspaceModuleSelectorView(workspace storage.Workspace) WorkspaceModuleSelectorView {
	selected := make(map[string]struct{})
	for _, m := range strings.Split(workspace.ActiveModules, ",") {
		if m != "" {
			selected[m] = struct{}{}
		}
	}

	return WorkspaceModuleSelectorView{
		workspace:      workspace,
		availableModules: module.GetAvailableModules(),
		selected:       selected,
	}
}

func (v WorkspaceModuleSelectorView) Init() tea.Cmd {
	return nil
}

type DoneWorkspaceModuleSelectorMsg struct {
	Workspace storage.Workspace
}

func (v WorkspaceModuleSelectorView) Update(msg tea.Msg) (WorkspaceModuleSelectorView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return v, func() tea.Msg {
				return DoneWorkspaceModuleSelectorMsg{Workspace: v.workspace}
			}
		case "up", "k":
			if v.cursor > 0 {
				v.cursor--
			}
		case "down", "j":
			if v.cursor < len(v.availableModules)-1 {
				v.cursor++
			}
		case "enter", " ":
			moduleName := v.availableModules[v.cursor]
			if _, ok := v.selected[moduleName]; ok {
				delete(v.selected, moduleName)
			} else {
				v.selected[moduleName] = struct{}{}
			}

			var activeModules []string
			for m := range v.selected {
				activeModules = append(activeModules, m)
			}
			v.workspace.ActiveModules = strings.Join(activeModules, ",")
		}
	}
	return v, nil
}

func (v WorkspaceModuleSelectorView) View() string {
	var s strings.Builder
	s.WriteString("Select active modules for this workspace (press space to toggle, enter to save):\n\n")

	for i, module := range v.availableModules {
		cursor := " "
		if v.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := v.selected[module]; ok {
			checked = "x"
		}

		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, cursor, " [", checked, "] ", module, "\n"))
	}

	return s.String()
}
