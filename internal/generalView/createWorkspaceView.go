package generalview

import (
	"database/sql"

	"github.com/Ceinl/Go-dashboard/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type CreateWorkspaceView struct {
	Width  int
	Height int

	db *sql.DB

	nameInput  textinput.Model
	colorInput textinput.Model
	focused    int // 0 for name, 1 for color, 2 for OK button
}

// A message to signal that the workspace creation is done/cancelled.
type DoneCreateWorkspaceMsg struct{}

func NewCreateWorkspaceView(db *sql.DB) CreateWorkspaceView {
	v := CreateWorkspaceView{
		db: db,
	}
	v.nameInput = textinput.New()
	v.nameInput.Placeholder = "Workspace Name"
	v.nameInput.Focus()
	v.nameInput.CharLimit = 20
	v.nameInput.Width = 20

	v.colorInput = textinput.New()
	v.colorInput.Placeholder = "Color (hex)"
	v.colorInput.CharLimit = 7
	v.colorInput.Width = 20

	return v
}

func (v CreateWorkspaceView) Init() tea.Cmd {
	return textinput.Blink
}

func (v CreateWorkspaceView) Update(msg tea.Msg) (CreateWorkspaceView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.Width = msg.Width + 2
		v.Height = msg.Height + 2
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

			if v.focused > 2 {
				v.focused = 0
			} else if v.focused < 0 {
				v.focused = 2
			}

			cmds = append(cmds, v.updateFocus())
			return v, tea.Batch(cmds...)
		case "enter":
			if v.focused == 2 {
				return v, func() tea.Msg {
					// Create new workspace
					newWorkspace := storage.Workspace{
						ID:    uuid.New().String(),
						Name:  v.nameInput.Value(),
						Color: v.colorInput.Value(),
					}
					err := storage.CreateWorkspace(v.db, newWorkspace)
					if err != nil {
						// TODO: Handle error more gracefully, e.g., show an error message to the user
						return DoneCreateWorkspaceMsg{}
					}
					return DoneCreateWorkspaceMsg{}
				}
			} else {
				v.focused++
				cmds = append(cmds, v.updateFocus())
				return v, tea.Batch(cmds...)
			}
		}
		// If not a special key, pass to the currently focused text input
		var inputCmd tea.Cmd
		if v.focused == 0 {
			v.nameInput, inputCmd = v.nameInput.Update(msg)
		} else if v.focused == 1 {
			v.colorInput, inputCmd = v.colorInput.Update(msg)
		}
		if inputCmd != nil {
			cmds = append(cmds, inputCmd)
		}
	}

	return v, tea.Batch(cmds...)
}

func (v *CreateWorkspaceView) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, 2)
	if v.focused == 0 {
		cmds[0] = v.nameInput.Focus()
		v.colorInput.Blur()
	} else if v.focused == 1 {
		cmds[0] = v.colorInput.Focus()
		v.nameInput.Blur()
	} else {
		v.nameInput.Blur()
		v.colorInput.Blur()
	}
	return tea.Batch(cmds...)
}

func (v CreateWorkspaceView) View() string {
	buttonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("7")).
		Padding(0, 3)

	if v.focused == 2 {
		buttonStyle = buttonStyle.Copy().
			Foreground(lipgloss.Color("7")).
			Background(lipgloss.Color("0"))
	}

	okButton := buttonStyle.Render("OK")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"Create New Workspace Wizard",
		v.nameInput.View(),
		"",
		v.colorInput.View(),
		"",
		okButton,
		"",
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(2, 4). // Increased padding for bigger wizard
		Render(content)

	return lipgloss.Place(v.Width, v.Height, lipgloss.Center, lipgloss.Center, box)
}
