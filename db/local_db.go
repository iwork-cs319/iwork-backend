package db

import "go-api/model"

type LocalDBStore struct {
	workspaces map[string]*model.Workspace
}

func (l LocalDBStore) GetOneWorkspace(id string) (*model.Workspace, error) {
	w, ok := l.workspaces[id]
	if !ok {
		return nil, NotFoundError
	}
	return w, nil
}

func (l LocalDBStore) UpdateWorkspace(id string, workspace *model.Workspace) error {
	_, ok := l.workspaces[id]
	if !ok {
		return NotFoundError
	}
	if workspace.Floor != "" {
		l.workspaces[id].Floor = workspace.Floor
	}
	if workspace.Name != "" {
		l.workspaces[id].Name = workspace.Name
	}
	if workspace.Props != nil {
		l.workspaces[id].Props = workspace.Props
	}
	return nil
}

func (l LocalDBStore) CreateWorkspace(workspace *model.Workspace) error {
	l.workspaces[workspace.ID] = workspace
	return nil
}

func (l LocalDBStore) RemoveWorkspace(id string) error {
	delete(l.workspaces, id)
	return nil
}

func (l LocalDBStore) GetAllWorkspaces() ([]*model.Workspace, error) {
	var list []*model.Workspace
	if len(l.workspaces) < 1 {
		return nil, EmptyError
	}
	for _, w := range l.workspaces {
		list = append(list, w)
	}
	return list, nil
}

func NewLocalDataStore() *DataStore {
	return &DataStore{
		WorkspaceProvider: &LocalDBStore{workspaces: map[string]*model.Workspace{
			"1": {
				ID:    "1",
				Name:  "Workspace 1",
				Props: nil,
			},
			"2": {
				ID:    "2",
				Name:  "Workspace 2",
				Props: nil,
			},
			"3": {
				ID:    "3",
				Name:  "Workspace 3",
				Props: nil,
			},
			"6": {
				ID:    "6",
				Name:  "Workspace 6",
				Props: nil,
			},
		}},
	}
}
