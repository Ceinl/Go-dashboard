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
	return fmt.Sprintf("Name: %s", i.workspace.Name)
}

type SwapWorkspaceView struct {
	list list.Model
	db   *sql.DB
}

type DoneSwapWorkspaceMsg struct {
	SelectedWorkspace storage.Workspace
}

func NewSwapWorkspaceView(db *sql.DB) SwapWorkspaceView {
	items := []list.Item{}
	workspaces, err := storage.GetAllWorkspaces(db)
	if err != nil {
		log.Printf("Error getting all workspaces for swap view: %v", err)
	} else {
		for _, ws := range workspaces {
			items = append(items, item{workspace: ws})
		}
	}

	delegate := list.NewDefaultDelegate()
	m := list.New(items, delegate, 20, 20)
	m.Title = "Select a Workspace"

	return SwapWorkspaceView{list: m, db: db}
}

func (v SwapWorkspaceView) Init() tea.Cmd {
	return nil
}

func (v SwapWorkspaceView) Update(msg tea.Msg) (SwapWorkspaceView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return v, func() tea.Msg { return DoneSwapWorkspaceMsg{} }
		case "enter":
			if i, ok := v.list.SelectedItem().(item); ok {
				return v, func() tea.Msg { return DoneSwapWorkspaceMsg{SelectedWorkspace: i.workspace} }
			}
		}
	}

	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return v, cmd
}

func (v SwapWorkspaceView) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(v.list.View())
}
