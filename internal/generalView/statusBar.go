package generalview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusBar struct {
	Width           int
	CommandMode     bool
	Command         string
	ActiveWorkspace string
	ActiveProject   string
}

// New message type for newWorkspace command
type NewWorkspaceCommandMsg struct{}
type DeleteWorkspaceCommandMsg struct{}
type SwapWorkspaceCommandMsg struct{}
type NewProjectCommandMsg struct{}
type SwapProjectCommandMsg struct{}
type HelpCommandMsg struct{}
type TwitterLoginCommandMsg struct{}
type TwitterPostCommandMsg struct{}
type DeleteProjectCommandMsg struct{}
type ModuleSelectorCommandMsg struct{}
type WorkspaceModuleSelectorCommandMsg struct{}

func (s StatusBar) Init() tea.Cmd {
	return nil
}

func (s StatusBar) Update(msg tea.Msg) (StatusBar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.Width = msg.Width
	case tea.KeyMsg:
		if s.CommandMode {
			switch msg.Type {
			case tea.KeyEnter:
				s.CommandMode = false
				cmd := s.Command
				s.Command = ""
				switch cmd {
				case "q", "Q", "quit", "Quit", "Exit", "exit":
					return s, tea.Quit
				case "newWorkspace", "neww":
					return s, func() tea.Msg { return NewWorkspaceCommandMsg{} }
				case "deleteWorkspace", "delw":
					return s, func() tea.Msg { return DeleteWorkspaceCommandMsg{} }
				case "swapWorkspace", "swapw":
					return s, func() tea.Msg { return SwapWorkspaceCommandMsg{} }
				case "newProject", "newp":
					return s, func() tea.Msg { return NewProjectCommandMsg{} }
				case "swapProject", "swapp":
					return s, func() tea.Msg { return SwapProjectCommandMsg{} }
				case "help":
					return s, func() tea.Msg { return HelpCommandMsg{} }
				case "login":
					return s, func() tea.Msg { return TwitterLoginCommandMsg{} }
				case "post":
					return s, func() tea.Msg { return TwitterPostCommandMsg{} }
				case "delp":
					return s, func() tea.Msg { return DeleteProjectCommandMsg{} }
				case "modules":
					return s, func() tea.Msg { return ModuleSelectorCommandMsg{} }
				case "config-modules":
					return s, func() tea.Msg { return WorkspaceModuleSelectorCommandMsg{} }
				}
			case tea.KeyEsc:
				s.CommandMode = false
				s.Command = ""
			case tea.KeyBackspace:
				if len(s.Command) > 0 {
					s.Command = s.Command[:len(s.Command)-1]
				}
			case tea.KeyRunes:
				s.Command += string(msg.Runes)
			}
		} else {
			if msg.String() == ":" {
				s.CommandMode = true
			}
			if msg.String() == "?" {
				return s, func() tea.Msg { return HelpCommandMsg{} }
			}
		}
	}
	return s, nil
}

func (s StatusBar) View() string {
	var left string
	if s.CommandMode {
		left = ":" + s.Command
	} else {
		left = "" // Don't show placeholder text
	}

	right := s.ActiveWorkspace + " > " + s.ActiveProject

	// Calculate available width for the left part
	leftWidth := s.Width - lipgloss.Width(right)

	// Style for the left part, taking up the remaining space
	leftStyle := lipgloss.NewStyle().
		PaddingLeft(1).
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("235")).
		Width(leftWidth)

	// Style for the right part
	rightStyle := lipgloss.NewStyle().
		PaddingRight(1).
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("235"))

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(left),
		rightStyle.Render(right),
	)
}