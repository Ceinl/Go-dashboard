package storage

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := InitDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	return db
}

func TestCreateAndGetAllWorkspaces(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a new workspace
	testWorkspace := Workspace{
		ID:        uuid.New().String(),
		Name:      "Test Workspace",
		Color:     "blue",
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	err := CreateWorkspace(db, testWorkspace)
	if err != nil {
		t.Fatalf("failed to create workspace: %v", err)
	}

	// Get all workspaces
	workspaces, err := GetAllWorkspaces(db)
	if err != nil {
		t.Fatalf("failed to get all workspaces: %v", err)
	}

	// Check if the created workspace is in the list
	if len(workspaces) != 1 {
		t.Fatalf("expected 1 workspace, got %d", len(workspaces))
	}

	if workspaces[0].Name != testWorkspace.Name {
		t.Errorf("expected workspace name %s, got %s", testWorkspace.Name, workspaces[0].Name)
	}

	// Clean up
	err = DeleteWorkspace(db, testWorkspace.ID)
	if err != nil {
		t.Fatalf("failed to delete workspace: %v", err)
	}
}

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}
