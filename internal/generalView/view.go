package generalview

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProjectBar struct {
	Width int
}

func (p ProjectBar) Init() tea.Cmd {
	return nil
}

func (p ProjectBar) Update(msg tea.Msg) (ProjectBar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.Width = msg.Width
	}
	return p, nil
}

func (p ProjectBar) View() string {
	content := "ProjectBar"
	return lipgloss.PlaceHorizontal(p.Width, lipgloss.Center, content)
}

// ─────────────────────────────────────────────

type StatusBar struct {
	Width           int
	CommandMode     bool
	Command         string
	ActiveWorkspace string
}

// New message type for newWorkspace command
type NewWorkspaceCommandMsg struct{}
type DeleteWorkspaceCommandMsg struct{}
type SwapWorkspaceCommandMsg struct{}

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
		}
	}
	return s, nil
}

func (s StatusBar) View() string {
	var content string
	if s.CommandMode {
		content = ":" + s.Command
	} else {
		content = "StatusBar"
	}

	activeWorkspaceInfo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("235")).
		PaddingRight(1).
		Render("Workspace: " + s.ActiveWorkspace)

	style := lipgloss.NewStyle().
		Width(s.Width).
		Align(lipgloss.Left).
		PaddingLeft(1).
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("235"))

	return lipgloss.JoinHorizontal(lipgloss.Bottom, style.Render(content), activeWorkspaceInfo)
}

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
