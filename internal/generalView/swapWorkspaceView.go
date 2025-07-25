package generalview

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Ceinl/Go-dashboard/internal/storage"
)

type item struct {
	workspace storage.Workspace
}

func (i item) FilterValue() string {
	return i.workspace.Name
}

func (i item) Title() string {
	return i.workspace.Name
}

func (i item) Description() string {
	return fmt.Sprintf("ID: %s, Color: %s, Created At: %s", i.workspace.ID, i.workspace.Color, i.workspace.CreatedAt)
}

type SwapWorkspaceView struct {
	Width  int
	Height int

	db   *sql.DB
	list list.Model
}

// A message to signal that the workspace swap is done/cancelled.
type DoneSwapWorkspaceMsg struct {
	SelectedWorkspace storage.Workspace
}

func NewSwapWorkspaceView(db *sql.DB) SwapWorkspaceView {
	v := SwapWorkspaceView{
		db: db,
	}

	items := []list.Item{}
	workspaces, err := storage.GetAllWorkspaces(db)
	if err != nil {
		log.Printf("Error getting all workspaces for swap view: %v", err)
	} else {
		for _, ws := range workspaces {
			items = append(items, item{workspace: ws})
		}
	}

	m := list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.Title = "Select a Workspace"
	m.SetShowStatusBar(false)
	m.SetFilteringEnabled(true)
	m.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Background(lipgloss.Color("235"))
	m.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	m.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)

	v.list = m

	return v
}

func (v SwapWorkspaceView) Init() tea.Cmd {
	return nil
}

func (v SwapWorkspaceView) Update(msg tea.Msg) (SwapWorkspaceView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.Width = msg.Width + 2
		v.Height = msg.Height + 2
		v.list.SetSize(msg.Width, msg.Height-2)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return v, func() tea.Msg { return DoneSwapWorkspaceMsg{} }
		case "enter":
			selItem, ok := v.list.SelectedItem().(item)
			if ok {
				return v, func() tea.Msg { return DoneSwapWorkspaceMsg{SelectedWorkspace: selItem.workspace} }
			}
		}
	}

	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	cmds = append(cmds, cmd)

	return v, tea.Batch(cmds...)
}

func (v SwapWorkspaceView) View() string {
	return lipgloss.Place(v.Width, v.Height, lipgloss.Center, lipgloss.Center, v.list.View())
}
