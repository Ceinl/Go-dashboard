package generalview

import (
	"database/sql"

	"github.com/Ceinl/Go-dashboard/internal/storage"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DoneDeleteWorkspaceMsg struct{}
type ConfirmDeleteWorkspaceMsg struct {
	WorkspaceID string
}

type DeleteWorkspaceView struct {
	Width  int
	Height int

	db   *sql.DB
	list list.Model
}

func NewDeleteWorkspaceView(db *sql.DB) DeleteWorkspaceView {
	v := DeleteWorkspaceView{
		db: db,
	}

	items := []list.Item{}
	workspaces, err := storage.GetAllWorkspaces(db)
	if err != nil {
		// TODO: Handle error
	} else {
		for _, ws := range workspaces {
			items = append(items, deleteItem{workspace: ws})
		}
	}

	m := list.New(items, list.NewDefaultDelegate(), 20, 10)
	m.Title = "Select a Workspace to Delete"
	m.SetShowStatusBar(false)
	m.SetFilteringEnabled(true)
	m.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Background(lipgloss.Color("235"))
	m.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	m.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)

	v.list = m

	return v
}

func (v DeleteWorkspaceView) Init() tea.Cmd {
	return nil
}

func (v DeleteWorkspaceView) Update(msg tea.Msg) (DeleteWorkspaceView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.Width = msg.Width
		v.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return v, func() tea.Msg { return DoneDeleteWorkspaceMsg{} }
		case "enter":
			selItem, ok := v.list.SelectedItem().(deleteItem)
			if ok {
				return v, func() tea.Msg {
					return ConfirmDeleteWorkspaceMsg{WorkspaceID: selItem.workspace.ID}
				}
			}
		}
	}

	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	cmds = append(cmds, cmd)

	return v, tea.Batch(cmds...)
}

func (v DeleteWorkspaceView) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"Delete Workspace",
		v.list.View(),
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(2, 4).
		Render(content)

	return lipgloss.Place(v.Width, v.Height, lipgloss.Center, lipgloss.Center, box)
}

type deleteItem struct {
	workspace storage.Workspace
}

func (i deleteItem) Title() string       { return i.workspace.Name }
func (i deleteItem) Description() string { return i.workspace.Color }
func (i deleteItem) FilterValue() string { return i.workspace.Name }