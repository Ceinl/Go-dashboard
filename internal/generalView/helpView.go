package generalview

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type helpItem struct {
	key         string
	description string
}

type HelpView struct {
	content string
}

func NewHelpView() HelpView {
	items := []helpItem{
		{key: ":neww", description: "Create a new workspace"},
		{key: ":newp", description: "Create a new project"},
		{key: ":delw", description: "Delete a workspace"},
		{key: ":delp", description: "Delete a project"},
		{key: ":swapw", description: "Swap a workspace"},
		{key: ":swapp", description: "Swap a project"},
		{key: ":help", description: "Show this help screen"},
		{key: ":config-modules", description: "Configure modules for a workspace"},
		{key: "shift+h/l", description: "Switch between projects"},
		{key: "ctrl+h/l", description: "Switch between modules"},
		{key: "ctrl+s", description: "Save tweet as draft"},
		{key: "a", description: "Add a new item (linksaver, kanban)"},
		{key: "d", description: "Delete an item (linksaver, kanban)"},
		{key: "p", description: "Paste from clipboard (linksaver)"},
		{key: "enter", description: "Open a link (linksaver)"},
	}

	var content strings.Builder
	content.WriteString("Help\n\n")
	for _, item := range items {
		content.WriteString(fmt.Sprintf("%-20s %s\n", item.key, item.description))
	}

	return HelpView{content: content.String()}
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
	return v, nil
}

func (v HelpView) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(v.content)
}

type DoneHelpMsg struct{}
