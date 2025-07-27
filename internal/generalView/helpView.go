package generalview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type helpItem struct {
	key         string
	description string
}

func (i helpItem) FilterValue() string { return i.key }
func (i helpItem) Title() string       { return i.key }
func (i helpItem) Description() string { return i.description }

type HelpView struct {
	list list.Model
}

func NewHelpView() HelpView {
	items := []list.Item{
		helpItem{key: ":neww", description: "Create a new workspace"},
		helpItem{key: ":newp", description: "Create a new project"},
		helpItem{key: ":modules", description: "Select active modules for a project"},
		helpItem{key: ":help or ?", description: "Show this help screen"},
		helpItem{key: "shift+h/l", description: "Switch between projects"},
		helpItem{key: "ctrl+h/l", description: "Switch between modules"},
		helpItem{key: "ctrl+t", description: "Focus tweet composer"},
		helpItem{key: "ctrl+d", description: "Focus drafts list (Twitter module)"},
		helpItem{key: "ctrl+h", description: "Focus timeline (Twitter module)"},
		helpItem{key: "ctrl+u", description: "Focus user timeline (Twitter module)"},
		helpItem{key: "ctrl+c", description: "Copy tweet to clipboard"},
		helpItem{key: "ctrl+s", description: "Save tweet as draft"},
		helpItem{key: "ctrl+p", description: "Post tweet"},
		helpItem{key: "a", description: "Add a new item (linksaver, kanban)"},
		helpItem{key: "c", description: "Create a new column (kanban)"},
		helpItem{key: "r", description: "Rename a column (kanban)"},
		helpItem{key: "d", description: "Delete an item (linksaver, kanban)"},
		helpItem{key: "dd", description: "Delete a task (kanban)"},
		helpItem{key: "ctrl+d", description: "Delete a column (kanban)"},
		helpItem{key: "e", description: "Edit a task (kanban)"},
		helpItem{key: "p", description: "Paste from clipboard (linksaver)"},
		helpItem{key: "enter", description: "Open a link (linksaver)"},
	}

	m := list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.Title = "Help"
	m.SetShowHelp(false)
	m.SetShowFilter(false)

	return HelpView{list: m}
}

func (v HelpView) Init() tea.Cmd {
	return nil
}

func (v HelpView) Update(msg tea.Msg) (HelpView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return v, func() tea.Msg { return DoneHelpMsg{} }
		}
	}

	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return v, cmd
}

func (v HelpView) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(v.list.View())
}

type DoneHelpMsg struct{}
