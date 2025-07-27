package generalview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmationView struct {
	message  string
	onYes    tea.Msg
	onNo     tea.Msg
	focused  bool // false for No, true for Yes
	Width    int
	Height   int
	quitting bool
}

func NewConfirmationView(message string, onYes, onNo tea.Msg) ConfirmationView {
	return ConfirmationView{
		message: message,
		onYes:   onYes,
		onNo:    onNo,
	}
}

func (v ConfirmationView) Init() tea.Cmd {
	return nil
}

func (v ConfirmationView) Update(msg tea.Msg) (ConfirmationView, tea.Cmd) {
	if v.quitting {
		return v, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.Width = msg.Width
		v.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "right", "l", "tab", "shift+tab":
			v.focused = !v.focused
		case "enter":
			v.quitting = true
			if v.focused {
				return v, func() tea.Msg { return v.onYes }
			}
			return v, func() tea.Msg { return v.onNo }
		case "esc":
			v.quitting = true
			return v, func() tea.Msg { return v.onNo }
		}
	}
	return v, nil
}

func (v ConfirmationView) View() string {
	yesButton := "Yes"
	noButton := "No"

	if v.focused {
		yesButton = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Render(yesButton)
	} else {
		noButton = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Render(noButton)
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, noButton, yesButton)
	question := lipgloss.NewStyle().Width(v.Width).Align(lipgloss.Center).Render(v.message)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	return lipgloss.Place(v.Width, v.Height, lipgloss.Center, lipgloss.Center, ui)
}
