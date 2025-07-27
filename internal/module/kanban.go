package module

import (
	"database/sql"
	"log"
	"strings"

	"github.com/Ceinl/Go-dashboard/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

const (
	ToDo       = "To Do"
	InProgress = "In Progress"
	Done       = "Done"
)

var columns = []string{ToDo, InProgress, Done}

type Kanban struct {
	db        *sql.DB
	projectID string
	tasks     map[string][]storage.Task
	input     textinput.Model
	editing   bool
	cursorCol int
	cursorRow int
	width     int
	height    int
}

func NewKanban(db *sql.DB, projectID string) Module {
	ti := textinput.New()
	ti.Placeholder = "New Task"
	ti.CharLimit = 256
	ti.Width = 50

	return &Kanban{
		db:        db,
		projectID: projectID,
		tasks:     make(map[string][]storage.Task),
		input:     ti,
	}
}

func (m *Kanban) Init() tea.Cmd {
	m.loadTasks()
	return nil
}

func (m *Kanban) Update(msg tea.Msg) (Module, tea.Cmd) {
	if m.editing {
		return m.updateEditing(msg)
	}

	return m.updateBrowsing(msg)
}

func (m *Kanban) updateEditing(msg tea.Msg) (Module, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.projectID != "" {
				newTask := storage.Task{
					ID:        uuid.New().String(),
					ProjectID: m.projectID,
					Title:     m.input.Value(),
					Status:    columns[m.cursorCol],
				}
				if err := storage.CreateTask(m.db, newTask); err != nil {
					log.Printf("Error creating task: %v", err)
				} else {
					m.tasks[newTask.Status] = append(m.tasks[newTask.Status], newTask)
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

func (m *Kanban) updateBrowsing(msg tea.Msg) (Module, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if m.cursorCol > 0 {
				m.cursorCol--
				m.cursorRow = 0
			}
		case "right", "l":
			if m.cursorCol < len(columns)-1 {
				m.cursorCol++
				m.cursorRow = 0
			}
		case "up", "k":
			if m.cursorRow > 0 {
				m.cursorRow--
			}
		case "down", "j":
			if m.cursorRow < len(m.tasks[columns[m.cursorCol]])-1 {
				m.cursorRow++
			}
		case "H":
			m.moveTask(-1)
		case "L":
			m.moveTask(1)
		case "a":
			m.editing = true
			m.input.Focus()
			return m, textinput.Blink
		case "d":
			m.deleteTask()
		}
	}
	return m, nil
}

func (m *Kanban) View() string {
	if m.width == 0 {
		return "loading..."
	}

	var colViews []string
	columnWidth := m.width / len(columns)

	for i, colName := range columns {
		var tasksInCol []string
		for j, task := range m.tasks[colName] {
			taskStyle := lipgloss.NewStyle().Padding(0, 1).Width(columnWidth - 4)
			if i == m.cursorCol && j == m.cursorRow {
				taskStyle = taskStyle.Background(lipgloss.Color("57"))
			}
			tasksInCol = append(tasksInCol, taskStyle.Render(task.Title))
		}

		colStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Width(columnWidth - 2).
			Height(m.height - 10)
		if i == m.cursorCol {
			colStyle = colStyle.BorderForeground(lipgloss.Color("57"))
		}

		colViews = append(colViews, colStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Render(colName),
				strings.Join(tasksInCol, "\n"),
			),
		))
	}

	if m.editing {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.input.View())
	}

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, colViews...)
	helpView := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render("\n(a)dd, (d)elete, (h/j/k/l) navigate, (H/L) move task")

	return lipgloss.JoinVertical(lipgloss.Left, mainView, helpView)
}

func (m *Kanban) loadTasks() {
	if m.projectID == "" {
		m.tasks = make(map[string][]storage.Task)
		return
	}
	tasks, err := storage.GetTasksForProject(m.db, m.projectID)
	if err != nil {
		log.Printf("Error loading tasks: %v", err)
		return
	}

	m.tasks = make(map[string][]storage.Task)
	for _, col := range columns {
		m.tasks[col] = []storage.Task{}
	}
	for _, task := range tasks {
		m.tasks[task.Status] = append(m.tasks[task.Status], task)
	}
}

func (m *Kanban) moveTask(direction int) {
	currentColName := columns[m.cursorCol]
	if len(m.tasks[currentColName]) == 0 {
		return
	}

	task := m.tasks[currentColName][m.cursorRow]
	newColIndex := m.cursorCol + direction
	if newColIndex < 0 || newColIndex >= len(columns) {
		return
	}

	// Remove from old column
	m.tasks[currentColName] = append(m.tasks[currentColName][:m.cursorRow], m.tasks[currentColName][m.cursorRow+1:]...)

	// Add to new column
	newColName := columns[newColIndex]
	task.Status = newColName
	m.tasks[newColName] = append(m.tasks[newColName], task)

	// Update in DB
	if err := storage.UpdateTask(m.db, task); err != nil {
		log.Printf("Error updating task: %v", err)
		// Revert if DB update fails
		m.loadTasks()
	}

	m.cursorCol = newColIndex
	m.cursorRow = len(m.tasks[newColName]) - 1
}

func (m *Kanban) deleteTask() {
	currentColName := columns[m.cursorCol]
	if len(m.tasks[currentColName]) == 0 {
		return
	}

	task := m.tasks[currentColName][m.cursorRow]
	if err := storage.DeleteTask(m.db, task.ID); err != nil {
		log.Printf("Error deleting task: %v", err)
	} else {
		m.tasks[currentColName] = append(m.tasks[currentColName][:m.cursorRow], m.tasks[currentColName][m.cursorRow+1:]...)
		if m.cursorRow >= len(m.tasks[currentColName]) && len(m.tasks[currentColName]) > 0 {
			m.cursorRow = len(m.tasks[currentColName]) - 1
		}
	}
}
