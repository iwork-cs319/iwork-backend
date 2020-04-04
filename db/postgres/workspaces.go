package postgres

import (
	"database/sql"
	"errors"
	"fmt"
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
	var count int
	existsStmt := `SELECT count(*) FROM workspaces WHERE name=$1 AND floor_id=$2`
	err = tx.QueryRow(existsStmt, workspace.Name, workspace.Floor).Scan(&count)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", errors.New(fmt.Sprintf("workspace name: %s already exists on floor: %s", workspace.Name, workspace.Floor))
	}
	createWorkspaceStmt :=
		`INSERT INTO workspaces(name, floor_id, metadata, details) VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRow(
		createWorkspaceStmt,
		workspace.Name,
		workspace.Floor,
		workspace.Props,
		workspace.Details,
	).Scan(&workspaceId)
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	return workspaceId, err
}

func (p PostgresDBStore) UpsertWorkspace(workspace *model.Workspace) (string, error) {
	var err error
	workspaceId := ""
	tx, err := p.database.Begin()
	defer tx.Rollback()
	if err != nil {
		return "", err
	}

	// Check if workspace exists
	if workspace.ID == "" {
		err = tx.QueryRow(
			`SELECT id from workspaces where name=$1 AND floor_id=$2`,
			workspace.Name,
			workspace.Floor,
		).Scan(&workspaceId)
		if err != nil && err != sql.ErrNoRows {
			// Error retrieving data
			return "", nil
		}
		if workspaceId != "" {
			workspace.ID = workspaceId
		}
	}

	if workspace.ID != "" {
		// Update
		createWorkspaceStmt :=
			`UPDATE workspaces SET name=$2, floor_id=$3, metadata=$4, details=$5 WHERE id=$1 RETURNING id`
		err = tx.QueryRow(
			createWorkspaceStmt,
			workspace.ID,
			workspace.Name,
			workspace.Floor,
			workspace.Props,
			workspace.Details,
		).Scan(&workspaceId)
		if err != nil {
			return "", err
		}
	} else {
		workspaceId, err = p.CreateWorkspace(workspace)
		if err != nil {
			return "", err
		}
	}
	return workspaceId, tx.Commit()
}

func (p PostgresDBStore) RemoveWorkspace(id string) error {
	panic("implement me")
}

func (p PostgresDBStore) GetAllWorkspaces() ([]*model.Workspace, error) {
	rows, err := p.database.Query(`SELECT id, name, floor_id, details, metadata FROM workspaces WHERE deleted=FALSE;`)
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
	rows, err := p.database.Query(`SELECT id, name, floor_id, details, metadata FROM workspaces WHERE floor_id=$1 AND deleted=FALSE;`, floorId)
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

func (p PostgresDBStore) CountWorkspacesByFloor(floorId string) (int, error) {
	sqlStatement := `SELECT COUNT(*) FROM workspaces WHERE floor_id=$1;`
	var numWorkspacesFloor int
	row := p.database.QueryRow(sqlStatement, floorId)
	err := row.Scan(&numWorkspacesFloor)
	if err != nil {
		return 0, err
	}
	return numWorkspacesFloor, nil
}

func (p PostgresDBStore) CreateAssignment(userId, workspaceId string) error {
	tx, err := p.database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	log.Printf("Postgres.CreateAssignment: time.Now()=%s", now.UTC().String())

	var count int
	if err = tx.QueryRow(
		`SELECT count(*) FROM workspace_assignee
					WHERE workspace_id=$1 AND user_id=$2 AND end_time IS NULL`,
		workspaceId, userId,
	).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Check for future bookings

	if err = tx.QueryRow(
		`SELECT count(*) FROM bookings 
						WHERE workspace_id=$1 AND end_time >= $2`,
		workspaceId, now,
	).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("invalid operation: workspace has outstanding bookings")
	}

	// End the assignments
	updateAssignmentsStmt := `UPDATE workspace_assignee SET end_time=$2 WHERE workspace_id=$1 AND end_time IS NULL RETURNING id`
	_, err = tx.Exec(updateAssignmentsStmt, workspaceId, now)
	if err != nil {
		log.Printf("PostgresDBStore.CreateAssignment: error updating older assignment: %v\n", err)
		return err
	}

	// End any default offerings
	updateDefaultOfferingsStmt := `UPDATE offerings SET end_time=$2 WHERE workspace_id=$1 AND end_time IS NULL RETURNING id`
	_, err = tx.Exec(updateDefaultOfferingsStmt, workspaceId, now)
	if err != nil {
		log.Printf("PostgresDBStore.CreateAssignment: error updating default future offerings: %v\n", err)
		return err
	}

	// update any non-default offerings (cancel them)
	updateOfferingsStmt := `UPDATE offerings SET cancelled=TRUE WHERE workspace_id=$1 AND end_time >= $2 RETURNING id`
	_, err = tx.Exec(updateOfferingsStmt, workspaceId, now)
	if err != nil {
		log.Printf("PostgresDBStore.CreateDefaultOffering: error updating future offerings: %v\n", err)
		return err
	}

	var id string
	sqlStatement := `INSERT INTO workspace_assignee(user_id, workspace_id, start_time) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(sqlStatement, userId, workspaceId, time.Now()).Scan(&id)
	if err != nil {
		log.Printf("PostgresDBStore.CreateDefaultOffering: error creating new assignment: %v\n", err)
		return err
	}
	return tx.Commit()
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
