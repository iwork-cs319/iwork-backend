package postgres

import (
	"database/sql"
	"errors"
	"go-api/model"
	"go-api/utils"
	"log"
	"time"
)

const BookingAdvanceTime = time.Hour * 24 * 30 * 6 // 6 months

func (p PostgresDBStore) GetOneWorkspace(id string) (*model.Workspace, error) {
	sqlStatement := `SELECT id, name, floor_id, details, metadata FROM workspaces WHERE id=$1;`
	var workspace model.Workspace
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(&workspace.ID, &workspace.Name, &workspace.Floor, &workspace.Details, &workspace.Props)
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (p PostgresDBStore) UpdateWorkspace(id string, workspace *model.Workspace) error {
	tx, err := p.database.Begin()
	defer tx.Rollback()
	if err != nil {
		return err
	}
	var count int
	existsStmt := `SELECT count(*) FROM workspaces WHERE name=$1 AND floor_id=$2 AND id <> $3`
	err = tx.QueryRow(existsStmt, workspace.Name, workspace.Floor, id).Scan(&count)
	if err != nil {
		return err
	}
	log.Println(count)
	log.Println(workspace)
	if count > 0 {
		return errors.New("workspace name already exists")
	}
	sqlStatement :=
		`UPDATE workspaces
				SET name = $2, floor_id = $3, details = $4
				WHERE id = $1
				RETURNING id, name, floor_id;`
	var _id string
	var name string
	var floorId string
	err = tx.QueryRow(sqlStatement, id, workspace.Name, workspace.Floor, workspace.Details).Scan(&_id, &name, &floorId)
	if err != nil {
		return err
	}
	if _id != id || name != workspace.Name || floorId != workspace.Floor {
		return CreateError
	}
	workspace.ID = _id
	return tx.Commit()
}

func (p PostgresDBStore) UpdateWorkspaceMetadata(id string, properties *model.Attrs) error {
	sqlStatement := `UPDATE workspaces SET metadata=$2 WHERE id=$1 RETURNING id`
	var _id string
	err := p.database.QueryRow(sqlStatement, id, properties).Scan(&_id)
	if err != nil {
		return err
	}
	if _id != id {
		return CreateError
	}
	return nil
}

func (p PostgresDBStore) CreateWorkspace(workspace *model.Workspace) (string, error) {
	tx, err := p.database.Begin()
	defer tx.Rollback()
	if err != nil {
		return "", err
	}
	var workspaceId string
	existsStmt := `SELECT id FROM workspaces WHERE name=$1 AND floor_id=$2`
	err = tx.QueryRow(existsStmt, workspace.Name, workspace.Floor).Scan(&workspaceId)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == sql.ErrNoRows {
		createWorkspaceStmt := `INSERT INTO workspaces(name, floor_id) VALUES ($1, $2) RETURNING id`
		err = p.database.QueryRow(createWorkspaceStmt, workspace.Name, workspace.Floor).Scan(&workspaceId)
		if err != nil {
			return "", err
		}
	}
	err = tx.Commit()
	return workspaceId, err
}

func (p PostgresDBStore) RemoveWorkspace(id string) error {
	panic("implement me")
}

func (p PostgresDBStore) GetAllWorkspaces() ([]*model.Workspace, error) {
	rows, err := p.database.Query(`SELECT id, name, floor_id, details, metadata FROM workspaces;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	workspaces := make([]*model.Workspace, 0)
	for rows.Next() {
		var workspace model.Workspace
		err = rows.Scan(&workspace.ID, &workspace.Name, &workspace.Floor, &workspace.Details, &workspace.Props)
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

func (p PostgresDBStore) GetAllWorkspacesByFloor(floorId string) ([]*model.Workspace, error) {
	rows, err := p.database.Query(`SELECT id, name, floor_id, details, metadata FROM workspaces WHERE floor_id=$1;`, floorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	workspaces := make([]*model.Workspace, 0)
	for rows.Next() {
		var workspace model.Workspace
		err = rows.Scan(&workspace.ID, &workspace.Name, &workspace.Floor, &workspace.Details, &workspace.Props)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.GetAllWorkspacesByFloorId: %v\n", err)
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

func (p PostgresDBStore) CreateAssignWorkspace(workspace *model.Workspace, userId string) (string, error) {
	tx, err := p.database.Begin()
	now := time.Now()
	cancelTime := now.Add(BookingAdvanceTime)
	existing := false
	defer tx.Rollback()
	if err != nil {
		return "", err
	}
	var workspaceId string
	existsStmt := `SELECT id FROM workspaces WHERE name=$1 AND floor_id=$2`
	err = tx.QueryRow(existsStmt, workspace.Name, workspace.Floor).Scan(&workspaceId)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == sql.ErrNoRows {
		// it doesnt exist; create it
		createWorkspaceStmt := `INSERT INTO workspaces(name, floor_id) VALUES ($1, $2) RETURNING id`
		err = tx.QueryRow(createWorkspaceStmt, workspace.Name, workspace.Floor).Scan(&workspaceId)
		if err != nil {
			return "", err
		}
		existing = false
	} else {
		// Cancel any default offerings
		var offeringId string
		var startTime time.Time
		err = tx.QueryRow(`SELECT id, start_time from offerings where workspace_id=$1 AND end_time IS NULL AND user_id=$2`, workspaceId, utils.EmptyUserUUID).Scan(&offeringId, &startTime)
		if err != nil && err != sql.ErrNoRows {
			return "", err
		}
		if startTime.After(cancelTime) {
			cancelTime = startTime.Add(24 * time.Hour)
		}
		if err != sql.ErrNoRows {
			updateStmt := `UPDATE offerings SET end_time=$2 WHERE id=$1 RETURNING id`
			var x string
			err = tx.QueryRow(updateStmt, offeringId, cancelTime).Scan(&x)
			if err != nil {
				log.Printf("PostgresDBStore.CreateAssignWorkspace: error updating existing default offering: %v\n", err)
			}
		}
		// Cancel any existing assignments
		var waId string
		err = tx.QueryRow(`SELECT id FROM workspace_assignee WHERE workspace_id=$1 AND end_time IS NULL`, workspaceId).Scan(&waId)
		if err != nil && err != sql.ErrNoRows {
			return "", err
		}
		if err != sql.ErrNoRows {
			updateStmt := `UPDATE workspace_assignee SET end_time=$2 WHERE id=$1 RETURNING id`
			var x string
			err = tx.QueryRow(updateStmt, waId, cancelTime).Scan(&x)
			if err != nil {
				log.Printf("PostgresDBStore.CreateAssignWorkspace: error updating older assignment: %v\n", err)
			}
		}
		existing = true
	}
	createTime := cancelTime.Add(time.Hour * 24)
	if !existing {
		createTime = now
	}
	if userId != "" {
		// Create an assignment
		var waId string
		createAssignmentStmt :=
			`INSERT INTO workspace_assignee(user_id, workspace_id, start_time) VALUES ($1, $2, $3) RETURNING id`
		err = tx.QueryRow(createAssignmentStmt, userId, workspaceId, createTime).Scan(&waId)
		if err != nil {
			return "", nil
		}
	} else {
		// Create an offering
		sqlStatement :=
			`INSERT INTO offerings(user_id, workspace_id, start_time, end_time, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id`
		var id string
		err := tx.QueryRow(sqlStatement, utils.EmptyUserUUID, workspaceId, createTime, nil, utils.EmptyUserUUID).Scan(&id)
		if err != nil {
			return "", err
		}
	}
	err = tx.Commit()
	return workspaceId, err
}
