package generalview

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCreateWorkspaceView_TextInput(t *testing.T) {
	v := NewCreateWorkspaceView()
	cmd := v.Init()
	if cmd != nil {
		cmd()
	}

	// Simulate typing "test name" into the name input
	msgs := []tea.KeyMsg{
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")},
	}

	for _, msg := range msgs {
		v, _ = v.Update(msg)
	}

	if v.nameInput.Value() != "test name" {
		t.Errorf("Expected name input to be 'test name', got %s", v.nameInput.Value())
	}

	// Simulate pressing Tab to switch to color input
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Simulate typing "#123456" into the color input
	msgs = []tea.KeyMsg{
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("#")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("1")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("3")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("4")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("5")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("6")},
	}

	for _, msg := range msgs {
		v, _ = v.Update(msg)
	}

	if v.colorInput.Value() != "#123456" {
		t.Errorf("Expected color input to be '#123456', got %s", v.colorInput.Value())
	}

	// Simulate pressing Tab to switch to OK button
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Simulate pressing Enter on OK button
	_, cmd = v.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Error("Expected a command to be returned when pressing Enter on OK button, got nil")
	}

	msg := cmd()
	if _, ok := msg.(DoneCreateWorkspaceMsg); !ok {
		t.Errorf("Expected DoneCreateWorkspaceMsg, got %T", msg)
	}
}
