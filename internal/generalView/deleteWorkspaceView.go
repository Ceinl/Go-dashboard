package generalview

import (
	"database/sql"
	"fmt"

	"github.com/Ceinl/Go-dashboard/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DeleteWorkspaceView struct {
	Width  int
	Height int

	db *sql.DB

	nameInput textinput.Model
	focused   int    // 0 for name, 1 for OK button
	status    string // To display messages like "Workspace deleted" or "Not found"
}

// A message to signal that the workspace deletion is done/cancelled.
type DoneDeleteWorkspaceMsg struct{}

func NewDeleteWorkspaceView(db *sql.DB) DeleteWorkspaceView {
	v := DeleteWorkspaceView{
		db: db,
	}
	v.nameInput = textinput.New()
	v.nameInput.Placeholder = "Workspace Name to Delete"
	v.nameInput.Focus()
	v.nameInput.CharLimit = 50
	v.nameInput.Width = 30

	return v
}

func (v DeleteWorkspaceView) Init() tea.Cmd {
	return textinput.Blink
}

func (v DeleteWorkspaceView) Update(msg tea.Msg) (DeleteWorkspaceView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.Width = msg.Width + 2
		v.Height = msg.Height + 2
		v.nameInput.Width = msg.Width / 3
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return v, func() tea.Msg { return DoneDeleteWorkspaceMsg{} }
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

			cmds = append(cmds, v.updateFocus())
			return v, tea.Batch(cmds...)
		case "enter":
			if v.focused == 1 {
				// Attempt to delete workspace
				workspaceName := v.nameInput.Value()
				if workspaceName == "" {
					v.status = "Please enter a workspace name."
					return v, nil
				}

				ws, err := storage.GetWorkspaceByName(v.db, workspaceName)
				if err != nil {
					if err == sql.ErrNoRows {
						v.status = fmt.Sprintf("Workspace '%s' not found.", workspaceName)
					} else {
						v.status = fmt.Sprintf("Error finding workspace: %v", err)
					}
					return v, nil
				}

				err = storage.DeleteWorkspace(v.db, ws.ID)
				if err != nil {
					v.status = fmt.Sprintf("Error deleting workspace: %v", err)
				} else {
					v.status = fmt.Sprintf("Workspace '%s' deleted successfully.", workspaceName)
					// Clear input after successful deletion
					v.nameInput.SetValue("")
				}
				return v, nil
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
		}
		if inputCmd != nil {
			cmds = append(cmds, inputCmd)
		}
	}

	return v, tea.Batch(cmds...)
}

func (v *DeleteWorkspaceView) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, 2)
	switch v.focused {
	case 0:
		cmds[0] = v.nameInput.Focus()
	case 1:
		v.nameInput.Blur()
	}
	return tea.Batch(cmds...)
}

func (v DeleteWorkspaceView) View() string {
	buttonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("7")).
		Padding(0, 3)

	if v.focused == 1 {
		buttonStyle = buttonStyle.Copy().
			Foreground(lipgloss.Color("7")).
			Background(lipgloss.Color("0"))
	}

	deleteButton := buttonStyle.Render("Delete")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"Delete Workspace",
		v.nameInput.View(),
		"",
		deleteButton,
		"",
		v.status,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(2, 4). // Increased padding for bigger wizard
		Render(content)

	return lipgloss.Place(v.Width, v.Height, lipgloss.Center, lipgloss.Center, box)
}
