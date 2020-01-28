package db

import "errors"

var NotFoundError = errors.New("not found")
var EmptyError = errors.New("empty")

type workspace struct {
	ID    string
	Name  string
	Props map[string]interface{}
}

type DataStore interface {
	workspaceProvider
}

type workspaceProvider interface {
	GetOneWorkspace(id string) (*workspace, error)
	UpdateWorkspace(*workspace) error
	CreateWorkspace(*workspace) error
	RemoveWorkspace(id string) error
	GetAllWorkspaces() ([]*workspace, error)
}
