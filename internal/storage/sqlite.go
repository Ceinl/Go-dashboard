package storage

import (
	"database/sql"

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
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT NOT NULL PRIMARY KEY,
		workspace_id TEXT NOT NULL,
		name TEXT,
		description TEXT,
		status TEXT,
		FOREIGN KEY(workspace_id) REFERENCES workspaces(id)
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type Workspace struct {
	ID        string
	Name      string
	Color     string
	CreatedAt string
}

func CreateWorkspace(db *sql.DB, workspace Workspace) error {
	stmt, err := db.Prepare("INSERT INTO workspaces(id, name, color) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(workspace.ID, workspace.Name, workspace.Color)
	return err
}

func GetWorkspace(db *sql.DB, id string) (Workspace, error) {
	row := db.QueryRow("SELECT id, name, color, created_at FROM workspaces WHERE id = ?", id)

	var workspace Workspace
	err := row.Scan(&workspace.ID, &workspace.Name, &workspace.Color, &workspace.CreatedAt)
	if err != nil {
		return Workspace{}, err
	}

	return workspace, nil
}

func UpdateWorkspace(db *sql.DB, workspace Workspace) error {
	stmt, err := db.Prepare("UPDATE workspaces SET name = ?, color = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(workspace.Name, workspace.Color, workspace.ID)
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
	rows, err := db.Query("SELECT id, name, color, created_at FROM workspaces")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []Workspace
	for rows.Next() {
		var workspace Workspace
		if err := rows.Scan(&workspace.ID, &workspace.Name, &workspace.Color, &workspace.CreatedAt); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}

func GetWorkspaceByName(db *sql.DB, name string) (Workspace, error) {
	row := db.QueryRow("SELECT id, name, color, created_at FROM workspaces WHERE name = ?", name)

	var workspace Workspace
	err := row.Scan(&workspace.ID, &workspace.Name, &workspace.Color, &workspace.CreatedAt)
	if err != nil {
		return Workspace{}, err
	}

	return workspace, nil
}
