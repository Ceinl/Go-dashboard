package module

import tea "github.com/charmbracelet/bubbletea"

// Module defines the interface for all modules.
type Module interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Module, tea.Cmd)
	View() string
}
