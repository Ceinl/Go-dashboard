// internal/workspace/workspace.go
package workspace

import (
	"github.com/Ceinl/Go-dashboard/internal/module"
	"github.com/google/uuid"
)

type ProjectStatus string

const (
	StatusActive     ProjectStatus = "Active"
	StatusInactive   ProjectStatus = "Inactive"
	StatusCompleated ProjectStatus = "Compleated"
	StatusDropped    ProjectStatus = "Dropped"
)

type Project struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      ProjectStatus `json:"status"`
}

type Workspace struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Projects []Project       `json:"projects"`
	Modules  []module.Module `json:"modules"`
}

func NewWorkspace(name string) *Workspace {
	return &Workspace{
		ID:       uuid.NewString(),
		Name:     name,
		Projects: make([]Project, 0),
		Modules:  make([]module.Module, 0),
	}
}

func NewProject(name, description string) Project {
	return Project{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Status:      StatusActive,
	}
}
