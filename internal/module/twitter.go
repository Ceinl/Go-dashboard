package module

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Ceinl/Go-dashboard/internal/storage"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

const maxTweetLength = 280

type Twitter struct {
	db         *sql.DB
	projectID  string
	drafts     list.Model
	editor     textarea.Model
	editing    bool
	isCreating bool
	width      int
	height     int
}

func NewTwitter(db *sql.DB, projectID string) Module {
	editor := textarea.New()
	editor.Placeholder = "Write your tweet..."
	editor.CharLimit = maxTweetLength

	drafts := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	drafts.Title = "Drafts"

	return &Twitter{
		db:         db,
		projectID:  projectID,
		editor:     editor,
		drafts:     drafts,
		isCreating: false,
	}
}

func (m *Twitter) Init() tea.Cmd {
	m.loadTweets()
	return nil
}

func (m *Twitter) Update(msg tea.Msg) (Module, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.drafts.SetWidth(m.width / 3)
		m.drafts.SetHeight(m.height - 5)
		m.editor.SetWidth(m.width * 2 / 3)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			if m.editing {
				m.saveDraft()
				m.editing = false
				m.isCreating = false
				m.editor.Reset()
			}
		case "enter":
			if !m.editing {
				m.editing = true
				m.isCreating = false
				selectedItem := m.drafts.SelectedItem()
				if selectedItem != nil {
					tweet := selectedItem.(storage.Tweet)
					m.editor.SetValue(tweet.Content)
				}
				return m, textarea.Blink
			}
		case "n":
			if !m.editing {
				m.editing = true
				m.isCreating = true
				m.editor.Reset()
				m.editor.Focus()
				return m, textarea.Blink
			}
		case "esc":
			if m.editing {
				m.editing = false
				m.isCreating = false
				m.editor.Reset()
			}
		}
	}

	if m.editing {
		m.editor, cmd = m.editor.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.drafts, cmd = m.drafts.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Twitter) View() string {
	if m.width == 0 {
		return "loading..."
	}

	var helpView string
	if m.editing {
		helpView = "(ctrl+s) save, (esc) cancel"
	} else {
		helpView = "(n)ew, (enter) edit, (j/k) navigate"
	}

	draftsView := m.drafts.View()
	editorView := m.editor.View() + fmt.Sprintf("\n\n%d/%d", len(m.editor.Value()), maxTweetLength)

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, draftsView, lipgloss.NewStyle().Width(m.width*2/3).Render(editorView))
	if !m.editing {
		mainView = lipgloss.JoinHorizontal(lipgloss.Top, draftsView, lipgloss.NewStyle().Width(m.width*2/3).Render("Select a draft to edit or press 'n' to create a new one."))
	}

	return lipgloss.JoinVertical(lipgloss.Left, mainView, lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(helpView))
}

func (m *Twitter) loadTweets() {
	if m.projectID == "" {
		m.drafts.SetItems([]list.Item{})
		return
	}
	tweets, err := storage.GetTweetsForProject(m.db, m.projectID)
	if err != nil {
		log.Printf("Error loading tweets: %v", err)
		return
	}

	items := make([]list.Item, len(tweets))
	for i, tweet := range tweets {
		items[i] = tweet
	}
	m.drafts.SetItems(items)
}

func (m *Twitter) saveDraft() {
	content := m.editor.Value()
	if m.projectID != "" && content != "" {
		if m.isCreating {
			// Create new draft
			newTweet := storage.Tweet{
				ID:        uuid.New().String(),
				ProjectID: m.projectID,
				Content:   content,
			}
			if err := storage.CreateTweet(m.db, newTweet); err != nil {
				log.Printf("Error creating tweet: %v", err)
			}
		} else {
			// Update existing draft
			selectedItem := m.drafts.SelectedItem()
			if selectedItem != nil {
				tweet := selectedItem.(storage.Tweet)
				tweet.Content = content
				if err := storage.UpdateTweet(m.db, tweet); err != nil {
					log.Printf("Error updating tweet: %v", err)
				}
			}
		}
		m.loadTweets()
	}
}
