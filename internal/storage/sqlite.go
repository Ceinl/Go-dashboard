package storage

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS workspaces (
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT,
		color TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		active_modules TEXT
	);
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT NOT NULL PRIMARY KEY,
		workspace_id TEXT NOT NULL,
		name TEXT,
		description TEXT,
		status TEXT,
		active_modules TEXT,
		FOREIGN KEY(workspace_id) REFERENCES workspaces(id)
	);
	CREATE TABLE IF NOT EXISTS links (
		id TEXT NOT NULL PRIMARY KEY,
		project_id TEXT NOT NULL,
		title TEXT,
		url TEXT,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT NOT NULL PRIMARY KEY,
		project_id TEXT NOT NULL,
		title TEXT,
		status TEXT,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);
	CREATE TABLE IF NOT EXISTS tweets (
		id TEXT NOT NULL PRIMARY KEY,
		project_id TEXT NOT NULL,
		content TEXT,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	// Check if the active_modules column exists in projects
	rows, err := db.Query("PRAGMA table_info(projects)")
	if err != nil {
		return err
	}
	defer rows.Close()

	var columnExists bool
	for rows.Next() {
		var cid, notnull, pk int
		var name, dtype, dflt_value sql.NullString
		if err := rows.Scan(&cid, &name, &dtype, &notnull, &dflt_value, &pk); err != nil {
			return err
		}
		if name.Valid && name.String == "active_modules" {
			columnExists = true
			break
		}
	}

	if !columnExists {
		log.Println("Running migration: adding active_modules to projects table")
		_, err := db.Exec("ALTER TABLE projects ADD COLUMN active_modules TEXT")
		if err != nil {
			return err
		}
	}

	// Check if the active_modules column exists in workspaces
	rows, err = db.Query("PRAGMA table_info(workspaces)")
	if err != nil {
		return err
	}
	defer rows.Close()

	columnExists = false
	for rows.Next() {
		var cid, notnull, pk int
		var name, dtype, dflt_value sql.NullString
		if err := rows.Scan(&cid, &name, &dtype, &notnull, &dflt_value, &pk); err != nil {
			return err
		}
		if name.Valid && name.String == "active_modules" {
			columnExists = true
			break
		}
	}

	if !columnExists {
		log.Println("Running migration: adding active_modules to workspaces table")
		_, err := db.Exec("ALTER TABLE workspaces ADD COLUMN active_modules TEXT")
		if err != nil {
			return err
		}
	}

	return nil
}

type Workspace struct {
	ID            string
	Name          string
	Color         string
	CreatedAt     string
	ActiveModules string
}

func CreateWorkspace(db *sql.DB, workspace Workspace) error {
	stmt, err := db.Prepare("INSERT INTO workspaces(id, name, color, active_modules) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(workspace.ID, workspace.Name, workspace.Color, workspace.ActiveModules)
	return err
}

func GetWorkspace(db *sql.DB, id string) (Workspace, error) {
	row := db.QueryRow("SELECT id, name, color, created_at, active_modules FROM workspaces WHERE id = ?", id)

	var workspace Workspace
	err := row.Scan(&workspace.ID, &workspace.Name, &workspace.Color, &workspace.CreatedAt, &workspace.ActiveModules)
	if err != nil {
		return Workspace{}, err
	}

	return workspace, nil
}

func UpdateWorkspace(db *sql.DB, workspace Workspace) error {
	stmt, err := db.Prepare("UPDATE workspaces SET name = ?, color = ?, active_modules = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(workspace.Name, workspace.Color, workspace.ActiveModules, workspace.ID)
	return err
}

func DeleteWorkspace(db *sql.DB, id string) error {
	stmt, err := db.Prepare("DELETE FROM workspaces WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func GetAllWorkspaces(db *sql.DB) ([]Workspace, error) {
	rows, err := db.Query("SELECT id, name, color, created_at, active_modules FROM workspaces")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []Workspace
	for rows.Next() {
		var workspace Workspace
		if err := rows.Scan(&workspace.ID, &workspace.Name, &workspace.Color, &workspace.CreatedAt, &workspace.ActiveModules); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}

func GetWorkspaceByName(db *sql.DB, name string) (Workspace, error) {
	row := db.QueryRow("SELECT id, name, color, created_at, active_modules FROM workspaces WHERE name = ?", name)

	var workspace Workspace
	err := row.Scan(&workspace.ID, &workspace.Name, &workspace.Color, &workspace.CreatedAt, &workspace.ActiveModules)
	if err != nil {
		return Workspace{}, err
	}

	return workspace, nil
}

// Project represents a project within a workspace
type Project struct {
	ID            string
	WorkspaceID   string
	Name          string
	Description   string
	Status        string
	ActiveModules string // Stored as a comma-separated string
}

// CreateProject adds a new project to the database
func CreateProject(db *sql.DB, project Project) error {
	stmt, err := db.Prepare("INSERT INTO projects(id, workspace_id, name, description, status, active_modules) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(project.ID, project.WorkspaceID, project.Name, project.Description, project.Status, project.ActiveModules)
	return err
}

// UpdateProject updates a project in the database
func UpdateProject(db *sql.DB, project Project) error {
	stmt, err := db.Prepare("UPDATE projects SET name = ?, description = ?, status = ?, active_modules = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(project.Name, project.Description, project.Status, project.ActiveModules, project.ID)
	return err
}

// GetAllProjectsForWorkspace retrieves all projects for a given workspace
func GetAllProjectsForWorkspace(db *sql.DB, workspaceID string) ([]Project, error) {
	rows, err := db.Query("SELECT id, workspace_id, name, description, status, active_modules FROM projects WHERE workspace_id = ?", workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.ID, &project.WorkspaceID, &project.Name, &project.Description, &project.Status, &project.ActiveModules); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

type Link struct {
	ID        string
	ProjectID string
	Title     string
	URL       string
}

func GetLinksForProject(db *sql.DB, projectID string) ([]Link, error) {
	rows, err := db.Query("SELECT id, project_id, title, url FROM links WHERE project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		if err := rows.Scan(&link.ID, &link.ProjectID, &link.Title, &link.URL); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func CreateLink(db *sql.DB, link Link) error {
	stmt, err := db.Prepare("INSERT INTO links(id, project_id, title, url) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(link.ID, link.ProjectID, link.Title, link.URL)
	return err
}

func DeleteLink(db *sql.DB, id string) error {
	stmt, err := db.Prepare("DELETE FROM links WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

type Task struct {
	ID        string
	ProjectID string
	Title     string
	Status    string
}

func GetTasksForProject(db *sql.DB, projectID string) ([]Task, error) {
	rows, err := db.Query("SELECT id, project_id, title, status FROM tasks WHERE project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.ProjectID, &task.Title, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func CreateTask(db *sql.DB, task Task) error {
	stmt, err := db.Prepare("INSERT INTO tasks(id, project_id, title, status) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.ID, task.ProjectID, task.Title, task.Status)
	return err
}

func UpdateTask(db *sql.DB, task Task) error {
	stmt, err := db.Prepare("UPDATE tasks SET title = ?, status = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Title, task.Status, task.ID)
	return err
}

func DeleteTask(db *sql.DB, id string) error {
	stmt, err := db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

type Tweet struct {
	ID        string
	ProjectID string
	Content   string
}

func (t Tweet) Title() string {
	lines := strings.Split(t.Content, "\n")
	if len(lines) > 0 {
		return lines[0]
	}
	return ""
}
func (t Tweet) Description() string { return "" }
func (t Tweet) FilterValue() string { return t.Content }

func GetTweetsForProject(db *sql.DB, projectID string) ([]Tweet, error) {
	rows, err := db.Query("SELECT id, project_id, content FROM tweets WHERE project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []Tweet
	for rows.Next() {
		var tweet Tweet
		if err := rows.Scan(&tweet.ID, &tweet.ProjectID, &tweet.Content); err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}

	return tweets, nil
}

func CreateTweet(db *sql.DB, tweet Tweet) error {
	stmt, err := db.Prepare("INSERT INTO tweets(id, project_id, content) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tweet.ID, tweet.ProjectID, tweet.Content)
	return err
}

func UpdateTweet(db *sql.DB, tweet Tweet) error {
	stmt, err := db.Prepare("UPDATE tweets SET content = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tweet.Content, tweet.ID)
	return err
}

func DeleteTweet(db *sql.DB, id string) error {
	stmt, err := db.Prepare("DELETE FROM tweets WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}