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
			items = append(items, item{workspace: ws})
		}
	}

	m := list.New(items, list.NewDefaultDelegate(), 0, 0)
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
		v.Width = msg.Width + 2
		v.Height = msg.Height + 2
		v.list.SetSize(msg.Width, msg.Height-2)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return v, func() tea.Msg { return DoneDeleteWorkspaceMsg{} }
		case "enter":
			selItem, ok := v.list.SelectedItem().(item)
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
	return lipgloss.Place(v.Width, v.Height, lipgloss.Center, lipgloss.Center, v.list.View())
}