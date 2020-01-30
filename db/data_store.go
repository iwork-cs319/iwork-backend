package db

import "errors"
import "go-api/model"

var NotFoundError = errors.New("not found")
var EmptyError = errors.New("empty")

type DataStore struct {
	WorkspaceProvider workspaceProvider
}

type workspaceProvider interface {
	GetOneWorkspace(id string) (*model.Workspace, error)
	UpdateWorkspace(id string, workspace *model.Workspace) error
	CreateWorkspace(workspace *model.Workspace) error
	RemoveWorkspace(id string) error
	GetAllWorkspaces() ([]*model.Workspace, error)
}
