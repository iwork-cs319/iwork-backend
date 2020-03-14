package postgres

import (
	"go-api/model"
	"log"
	"time"
)

func (p PostgresDBStore) GetOneWorkspace(id string) (*model.Workspace, error) {
	sqlStatement := `SELECT id, name, floor_id FROM workspaces WHERE id=$1;`
	var workspace model.Workspace
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(&workspace.ID, &workspace.Name, &workspace.Floor)
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (p PostgresDBStore) UpdateWorkspace(id string, workspace *model.Workspace) error {
	sqlStatement :=
		`UPDATE workspaces
				SET name = $2, floor_id = $3
				WHERE id = $1
				RETURNING id, name, floor_id;`
	var _id string
	var name string
	var floorId string
	err := p.database.QueryRow(sqlStatement, id, workspace.Name, workspace.Floor).Scan(&_id, &name, &floorId)
	if err != nil {
		return err
	}
	if _id != id || name != workspace.Name || floorId != workspace.Floor {
		return CreateError
	}
	workspace.ID = _id
	return nil
}

func (p PostgresDBStore) CreateWorkspace(workspace *model.Workspace) (string, error) {
	sqlStatement := `INSERT INTO workspaces(name, floor_id) VALUES ($1, $2) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement, workspace.Name, workspace.Floor).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p PostgresDBStore) RemoveWorkspace(id string) error {
	panic("implement me")
}

func (p PostgresDBStore) GetAllWorkspaces() ([]*model.Workspace, error) {
	rows, err := p.database.Query(`SELECT id, name, floor_id FROM workspaces;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	workspaces := make([]*model.Workspace, 0)
	for rows.Next() {
		var workspace model.Workspace
		err = rows.Scan(&workspace.ID, &workspace.Name, &workspace.Floor)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.GetAllWorkspaces: %v\n", err)
		}
		workspaces = append(workspaces, &workspace)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (p PostgresDBStore) CreateAssignment(userId, workspaceId string) error {
	rows, err := p.database.Query(`SELECT id FROM workspace_assignee WHERE workspace_id=$1 AND end_time IS NOT NULL`, workspaceId)
	if err != nil {
		log.Printf("PostgresDBStore.CreateAssignment: error fetching older assignment: %v\n", err)
	}
	defer rows.Close()
	notEmpty := rows.Next()
	var id string
	if !notEmpty {
		err = rows.Scan(&id)
	}
	if id != "" {
		updateStmt :=
			`UPDATE workspace_assignee 
				SET end_time=$2 WHERE id=$1`
		err = p.database.QueryRow(updateStmt, id, time.Now()).Scan(&id)
		if err != nil {
			log.Printf("PostgresDBStore.CreateAssignment: error updating older assignment: %v\n", err)
		}
	}
	sqlStatement := `INSERT INTO workspace_assignee(user_id, workspace_id, start_time) VALUES ($1, $2, $3) RETURNING id`
	err = p.database.QueryRow(sqlStatement, userId, workspaceId, time.Now()).Scan(&id)
	return err
}
