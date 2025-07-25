package generalview

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
)

type CreateWorkspaceView struct {
	Width  int
	Height int

	nameInput  textinput.Model
	colorInput textinput.Model
	focused    int // 0 for name, 1 for color
}

// A message to signal that the workspace creation is done/cancelled.
type DoneCreateWorkspaceMsg struct{}

func (v CreateWorkspaceView) Init() tea.Cmd {
	v.nameInput = textinput.New()
	v.nameInput.Placeholder = "Workspace Name"
	v.nameInput.Focus()
	v.nameInput.CharLimit = 20
	v.nameInput.Width = 20

	v.colorInput = textinput.New()
	v.colorInput.Placeholder = "Workspace Color (e.g., #RRGGBB)"
	v.colorInput.CharLimit = 7
	v.colorInput.Width = 20

	return textinput.Blink
}

func (v CreateWorkspaceView) Update(msg tea.Msg) (CreateWorkspaceView, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.Width = msg.Width
		v.Height = msg.Height
		v.nameInput.Width = msg.Width / 3
		v.colorInput.Width = msg.Width / 3
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return v, func() tea.Msg { return DoneCreateWorkspaceMsg{} }
		case "tab", "shift+tab":
			if msg.String() == "shift+tab" {
				v.focused--
			} else {
				v.focused++
			}
			if v.focused > 1 {
				v.focused = 0
			} else if v.focused < 0 {
				v.focused = 1
			}

			if v.focused == 0 {
				cmds = append(cmds, v.nameInput.Focus())
				v.colorInput.Blur()
			} else {
				cmds = append(cmds, v.colorInput.Focus())
				v.nameInput.Blur()
			}
			return v, tea.Batch(cmds...)
		case "enter":
			if v.focused == 1 {
				return v, func() tea.Msg { return DoneCreateWorkspaceMsg{} }
			} else {
				v.focused = 1
				v.nameInput.Blur()
				cmds = append(cmds, v.colorInput.Focus())
				return v, tea.Batch(cmds...)
			}
		}
		// If not a special key, pass to the currently focused text input
		if v.focused == 0 {
			v.nameInput, cmd = v.nameInput.Update(msg)
		} else {
			v.colorInput, cmd = v.colorInput.Update(msg)
		}
		cmds = append(cmds, cmd)
	}

	return v, tea.Batch(cmds...)
}

func (v CreateWorkspaceView) View() string {
	content := fmt.Sprintf(
		"Create New Workspace Wizard\n\n%s\n%s\n\n(Press 'tab' to switch fields, 'enter' to submit, 'q' or 'esc' to close)",
		v.nameInput.View(),
		v.colorInput.View(),
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(v.Width, v.Height, lipgloss.Center, lipgloss.Center, box)
}
