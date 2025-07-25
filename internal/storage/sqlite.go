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
		name TEXT
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
