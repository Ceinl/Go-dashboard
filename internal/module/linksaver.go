package module

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/Ceinl/Go-dashboard/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type LinkSaver struct {
	db        *sql.DB
	projectID string
	links     []storage.Link
	input     textinput.Model
	editing   bool
	cursor    int
}

func NewLinkSaver(db *sql.DB, projectID string) Module {
	ti := textinput.New()
	ti.Placeholder = "Title,URL"
	ti.CharLimit = 256
	ti.Width = 50

	return &LinkSaver{
		db:        db,
		projectID: projectID,
		input:     ti,
	}
}

func (m *LinkSaver) Init() tea.Cmd {
	m.loadLinks()
	return nil
}

func (m *LinkSaver) Update(msg tea.Msg) (Module, tea.Cmd) {
	if m.editing {
		return m.updateEditing(msg)
	}

	return m.updateBrowsing(msg)
}

func (m *LinkSaver) updateEditing(msg tea.Msg) (Module, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.projectID != "" {
				linkData := strings.Split(m.input.Value(), ",")
				if len(linkData) == 2 {
					newLink := storage.Link{
						ID:        uuid.New().String(),
						ProjectID: m.projectID,
						Title:     strings.TrimSpace(linkData[0]),
						URL:       strings.TrimSpace(linkData[1]),
					}
					if err := storage.CreateLink(m.db, newLink); err != nil {
						log.Printf("Error creating link: %v", err)
					} else {
						m.links = append(m.links, newLink)
					}
				}
			}
			m.input.Reset()
			m.editing = false
			return m, nil
		case "esc":
			m.input.Reset()
			m.editing = false
			return m, nil
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *LinkSaver) updateBrowsing(msg tea.Msg) (Module, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.links)-1 {
				m.cursor++
			}
		case "a":
			m.editing = true
			m.input.Focus()
			return m, textinput.Blink
		case "p":
			clipboardContent, err := exec.Command("pbpaste").Output()
			if err != nil {
				log.Printf("Error pasting from clipboard: %v", err)
				return m, nil
			}
			m.input.SetValue("," + string(clipboardContent))
			m.editing = true
			m.input.Focus()
			return m, textinput.Blink
		case "d":
			if len(m.links) > 0 && m.cursor < len(m.links) {
				linkToDelete := m.links[m.cursor]
				if err := storage.DeleteLink(m.db, linkToDelete.ID); err != nil {
					log.Printf("Error deleting link: %v", err)
				} else {
					m.links = append(m.links[:m.cursor], m.links[m.cursor+1:]...)
					if m.cursor >= len(m.links) && len(m.links) > 0 {
						m.cursor = len(m.links) - 1
					}
				}
			}
		case "enter":
			if len(m.links) > 0 && m.cursor < len(m.links) {
				linkToOpen := m.links[m.cursor]
				exec.Command("open", linkToOpen.URL).Start()
			}
		case "c":
			if len(m.links) > 0 && m.cursor < len(m.links) {
				linkToCopy := m.links[m.cursor]
				cmd := exec.Command("pbcopy")
				cmd.Stdin = strings.NewReader(linkToCopy.URL)
				cmd.Run()
			}
		}
	}
	return m, nil
}

func (m *LinkSaver) View() string {
	var s strings.Builder
	s.WriteString("Link Saver\n\n")

	for i, link := range m.links {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s: %s\n", cursor, link.Title, link.URL))
	}

	if m.editing {
		s.WriteString("\n" + m.input.View())
	}

	s.WriteString("\n\n(a)dd, (p)aste, (d)elete, (c)opy, (enter) open, (j/k) navigate")
	return s.String()
}

func (m *LinkSaver) loadLinks() {
	if m.projectID == "" {
		m.links = []storage.Link{}
		return
	}
	links, err := storage.GetLinksForProject(m.db, m.projectID)
	if err != nil {
		log.Printf("Error loading links: %v", err)
		return
	}
	m.links = links
}
