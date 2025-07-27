package generalview

import (
	"fmt"
	"io"
	"strings"

	"github.com/Ceinl/Go-dashboard/internal/storage"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var availableModules = []string{"linksaver", "kanban", "twitter"}

type moduleItem struct {
	name     string
	selected bool
}

func (i moduleItem) FilterValue() string { return i.name }

type moduleDelegate struct{}

func (d moduleDelegate) Height() int                               { return 1 }
func (d moduleDelegate) Spacing() int                              { return 0 }
func (d moduleDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d moduleDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(moduleItem)
	if !ok {
		return
	}

	checkbox := "[ ]"
	if i.selected {
		checkbox = "[x]"
	}

	str := fmt.Sprintf("%s %s", checkbox, i.name)

	fn := func(s ...string) string {
		return lipgloss.NewStyle().PaddingLeft(4).Render(s...)
	}
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Render("> " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}

type ModuleSelectorView struct {
	list    list.Model
	project storage.Project
}

type DoneModuleSelectorMsg struct {
	Project storage.Project
}

func NewModuleSelectorView(project storage.Project) ModuleSelectorView {
	selectedModules := strings.Split(project.ActiveModules, ",")
	items := make([]list.Item, len(availableModules))
	for i, modName := range availableModules {
		item := moduleItem{name: modName}
		for _, sm := range selectedModules {
			if sm == modName {
				item.selected = true
			}
		}
		items[i] = item
	}

	m := list.New(items, moduleDelegate{}, 20, 20)
	m.Title = "Select Modules"
	m.SetShowHelp(false)

	return ModuleSelectorView{list: m, project: project}
}

func (v ModuleSelectorView) Init() tea.Cmd {
	return nil
}

func (v ModuleSelectorView) Update(msg tea.Msg) (ModuleSelectorView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			// Save and exit
			var selected []string
			for _, item := range v.list.Items() {
				if item.(moduleItem).selected {
					selected = append(selected, item.(moduleItem).name)
				}
			}
			v.project.ActiveModules = strings.Join(selected, ",")
			return v, func() tea.Msg { return DoneModuleSelectorMsg{Project: v.project} }
		case "enter", " ":
			// Toggle selection
			if i, ok := v.list.SelectedItem().(moduleItem); ok {
				i.selected = !i.selected
				v.list.SetItem(v.list.Index(), i)
			}
		}
	}

	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return v, cmd
}

func (v ModuleSelectorView) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(v.list.View())
}
