package postgres

import (
	"database/sql"
	"go-api/model"
	"log"
)

func (p PostgresDBStore) GetOneWorkspace(id string) (*model.Workspace, error) {
	sqlStatement := `SELECT id, name, floor_id, user_id FROM workspaces WHERE id=$1;`
	var workspace model.Workspace
	var userId sql.NullString
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(&workspace.ID, &workspace.Name, &workspace.Floor, &userId)
	if err != nil {
		return nil, err
	}
	if userId.Valid {
		workspace.User = userId.String
	}
	return &workspace, nil
}

func (p PostgresDBStore) UpdateWorkspace(id string, workspace *model.Workspace) error {
	sqlStatement :=
		`UPDATE workspaces
				SET name = $2, floor_id = $3
				WHERE id = $1
				RETURNING id, name;`
	var _id string
	var name string
	err := p.database.QueryRow(sqlStatement, id, workspace.Name, workspace.Floor).Scan(&_id, &name)
	if err != nil {
		return err
	}
	if _id != id || name != workspace.Name {
		return CreateError
	}
	return nil
}

func (p PostgresDBStore) CreateWorkspace(workspace *model.Workspace) error {
	sqlStatement := `INSERT INTO workspaces(id, name, floor_id) VALUES ($1, $2, $3) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement, workspace.ID, workspace.Name, workspace.Floor).Scan(&id)
	if err != nil {
		return err
	}
	if id != workspace.ID {
		return CreateError
	}
	return nil
}

func (p PostgresDBStore) RemoveWorkspace(id string) error {
	panic("implement me")
}

func (p PostgresDBStore) GetAllWorkspaces() ([]*model.Workspace, error) {
	rows, err := p.database.Query(`SELECT id, name, floor_id, user_id FROM workspaces;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	workspaces := make([]*model.Workspace, 0)
	for rows.Next() {
		var workspace model.Workspace
		var userId sql.NullString
		err = rows.Scan(&workspace.ID, &workspace.Name, &workspace.Floor, &userId)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.GetAllWorkspaces: %v\n", err)
		}
		if userId.Valid {
			workspace.User = userId.String
		}
		workspaces = append(workspaces, &workspace)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}
