package generalview

import (
	"database/sql"
	"log"

	"github.com/Ceinl/Go-dashboard/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type CreateProjectView struct {
	Width  int
	Height int

	db          *sql.DB
	workspaceID string

	nameInput  textinput.Model
	descInput  textinput.Model
	focused    int // 0 for name, 1 for desc, 2 for OK button
}

// A message to signal that the project creation is done/cancelled.
type DoneCreateProjectMsg struct{}

func NewCreateProjectView(db *sql.DB, workspaceID string) CreateProjectView {
	v := CreateProjectView{
		db:          db,
		workspaceID: workspaceID,
	}
	v.nameInput = textinput.New()
	v.nameInput.Placeholder = "Project Name"
	v.nameInput.Focus()
	v.nameInput.CharLimit = 30
	v.nameInput.Width = 30

	v.descInput = textinput.New()
	v.descInput.Placeholder = "Description"
	v.descInput.CharLimit = 100
	v.descInput.Width = 30

	return v
}

func (v CreateProjectView) Init() tea.Cmd {
	return textinput.Blink
}

func (v CreateProjectView) Update(msg tea.Msg) (CreateProjectView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.Width = msg.Width + 2
		v.Height = msg.Height + 2
		v.nameInput.Width = msg.Width / 3
		v.descInput.Width = msg.Width / 3
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return v, func() tea.Msg { return DoneCreateProjectMsg{} }
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
					newProject := storage.Project{
						ID:          uuid.New().String(),
						WorkspaceID: v.workspaceID,
						Name:        v.nameInput.Value(),
						Description: v.descInput.Value(),
					}
					err := storage.CreateProject(v.db, newProject)
					if err != nil {
						log.Printf("Error creating project: %v", err)
						// In a real app, you'd return an error message to the user
					}
					return DoneCreateProjectMsg{}
				}
			} else {
				v.focused++
				cmds = append(cmds, v.updateFocus())
				return v, tea.Batch(cmds...)
			}
		}
		var inputCmd tea.Cmd
		if v.focused == 0 {
			v.nameInput, inputCmd = v.nameInput.Update(msg)
		} else if v.focused == 1 {
			v.descInput, inputCmd = v.descInput.Update(msg)
		}
		if inputCmd != nil {
			cmds = append(cmds, inputCmd)
		}
	}

	return v, tea.Batch(cmds...)
}

func (v *CreateProjectView) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, 2)
	if v.focused == 0 {
		cmds[0] = v.nameInput.Focus()
		v.descInput.Blur()
	} else if v.focused == 1 {
		cmds[0] = v.descInput.Focus()
		v.nameInput.Blur()
	} else {
		v.nameInput.Blur()
		v.descInput.Blur()
	}
	return tea.Batch(cmds...)
}

func (v CreateProjectView) View() string {
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
		"Create New Project",
		v.nameInput.View(),
		"",
		v.descInput.View(),
		"",
		okButton,
		"",
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(2, 4).
		Render(content)

	return lipgloss.Place(v.Width, v.Height, lipgloss.Center, lipgloss.Center, box)
}