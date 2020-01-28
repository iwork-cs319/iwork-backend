package db

type LocalDBStore struct {
	workspaces map[string]*workspace
}

func (l LocalDBStore) GetOneWorkspace(id string) (*workspace, error) {
	w, ok := l.workspaces[id]
	if !ok {
		return nil, NotFoundError
	}
	return w, nil
}

func (l LocalDBStore) UpdateWorkspace(workspace *workspace) error {
	_, ok := l.workspaces[workspace.ID]
	if !ok {
		return NotFoundError
	}
	l.workspaces[workspace.ID] = workspace
	return nil
}

func (l LocalDBStore) CreateWorkspace(workspace *workspace) error {
	l.workspaces[workspace.ID] = workspace
	return nil
}

func (l LocalDBStore) RemoveWorkspace(id string) error {
	delete(l.workspaces, id)
	return nil
}

func (l LocalDBStore) GetAllWorkspaces() ([]*workspace, error) {
	var list []*workspace
	if len(l.workspaces) < 1 {
		return nil, EmptyError
	}
	for _, w := range l.workspaces {
		list = append(list, w)
	}
	return list, nil
}

func NewLocalDataStore() DataStore {
	return &LocalDBStore{workspaces: map[string]*workspace{
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
	}}
}
