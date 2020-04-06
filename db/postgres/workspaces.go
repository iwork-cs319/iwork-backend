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

	var assignedWorkspaceId string
	err = tx.QueryRow(
		`SELECT workspace_id FROM workspace_assignee
					WHERE user_id=$1 AND end_time IS NULL`,
		userId,
	).Scan(&assignedWorkspaceId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err != sql.ErrNoRows { // some assignment exists for this user
		if assignedWorkspaceId == workspaceId {
			// Already assigned to this workspace
			return nil
		}
		if assignedWorkspaceId != "" {
			// Already assigned to a workspace
			return nil
		}
	}

	// Check for future bookings
	var count int
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
	// Create workspace if doesnt exist; new=true if newly created
	// If userId attached is empty;
	// 		-> if new create default offering
	// 		-> else
	//			-> if currently assigned -> cancel assignment and create default offering
	//			-> else do nothing
	// Else (create new assignment)
	// 		-> if assigned to userID do nothing
	// 		-> if new create assignment
	// 		-> else (change of assignment)
	// 			-> if no future bookings -> cancel current assignment and offering; create assignment
	// 			-> else fail
	tx, err := p.database.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	now := time.Now()
	log.Printf("Postgres.CreateAssignWorkspace: time.Now()=%s", now.UTC().String())
	log.Printf("--- u:%s, n:%s\n", userId, workspace.Name)
	var workspaceId string
	newWorkspaceCreated := false
	existsStmt := `SELECT id FROM workspaces WHERE name=$1 AND floor_id=$2`
	err = tx.QueryRow(existsStmt, workspace.Name, workspace.Floor).Scan(&workspaceId)
	if err != nil && err != sql.ErrNoRows {
		// something bad happened; during query
		return "", err
	}
	if err == sql.ErrNoRows {
		// workspace doesnt exist; create it
		createWorkspaceStmt := `INSERT INTO workspaces(name, floor_id) VALUES ($1, $2) RETURNING id`
		err = tx.QueryRow(createWorkspaceStmt, workspace.Name, workspace.Floor).Scan(&workspaceId)
		if err != nil {
			return "", err
		}
		newWorkspaceCreated = true
		workspace.ID = workspaceId
		log.Println("--- Created new workspace", workspace)
	}

	currentlyAssignedUserId := ""
	if !newWorkspaceCreated {
		err = tx.QueryRow(`SELECT user_id FROM workspace_assignee WHERE workspace_id=$1 AND end_time IS NULL`, workspaceId).Scan(&currentlyAssignedUserId)
		if err != nil && err != sql.ErrNoRows {
			return "", err
		}
		if err == sql.ErrNoRows {
			currentlyAssignedUserId = ""
		}
	}
	log.Println("--- Currently Assigned userID: ", currentlyAssignedUserId)
	if userId == "" { // Create offering
		if newWorkspaceCreated {
			// create default offering
			sqlStatement :=
				`INSERT INTO offerings(user_id, workspace_id, start_time, end_time, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id`
			var id string
			err := tx.QueryRow(sqlStatement, utils.EmptyUserUUID, workspaceId, now, nil, utils.EmptyUserUUID).Scan(&id)
			if err != nil {
				return "", err
			}
			log.Println("--- Created New Default offering: ", id)
		} else {
			if currentlyAssignedUserId != "" {
				// cancel assignment
				updateStmt :=
					`UPDATE workspace_assignee SET end_time=$2 WHERE end_time IS NULL AND workspace_id=$1 AND user_id=$3 RETURNING id`
				var x string
				err = tx.QueryRow(updateStmt, workspaceId, now, currentlyAssignedUserId).Scan(&x)
				if err != nil {
					log.Printf("PostgresDBStore.CreateAssignWorkspace: error updating older assignment: %v\n", err)
					return "", err
				}
				log.Println("--- Cancelled assignment: ", x)
				// create default offering
				sqlStatement :=
					`INSERT INTO offerings(user_id, workspace_id, start_time, end_time, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id`
				var id string
				err := tx.QueryRow(sqlStatement, utils.EmptyUserUUID, workspaceId, now, nil, utils.EmptyUserUUID).Scan(&id)
				if err != nil {
					return "", err
				}
				log.Println("--- Created New Default offering: ", id)
			}
			// already offered -> do nothing
		}
	} else { // Create new assignment
		if newWorkspaceCreated {
			// Create assignment
			var waId string
			createAssignmentStmt :=
				`INSERT INTO workspace_assignee(user_id, workspace_id, start_time) VALUES ($1, $2, $3) RETURNING id`
			log.Println("=== ", userId, workspaceId, now)
			err = tx.QueryRow(createAssignmentStmt, userId, workspaceId, now).Scan(&waId)
			if err != nil {
				return "", err
			}
			log.Println("--- Created Assignment: ", waId)
		} else if currentlyAssignedUserId != userId {
			// check for future bookings on this workspace;
			var count int
			if err = tx.QueryRow(
				`SELECT count(*) FROM bookings 
						WHERE workspace_id=$1 AND end_time >= $2`,
				workspaceId, now,
			).Scan(&count); err != nil {
				return "", err
			}
			if count > 0 {
				return "", errors.New("invalid operation: workspace has outstanding bookings")
			}
			// cancel current assignment and offering
			updateAssignmentsStmt := `UPDATE workspace_assignee SET end_time=$2 WHERE workspace_id=$1 AND end_time IS NULL RETURNING id`
			_, err = tx.Exec(updateAssignmentsStmt, workspaceId, now)
			if err != nil {
				log.Printf("PostgresDBStore.CreateAssignment: error updating older assignment: %v\n", err)
				return "", err
			}
			log.Println("--- Cancelled all assignments ")

			// End any default offerings
			updateDefaultOfferingsStmt := `UPDATE offerings SET end_time=$2 WHERE workspace_id=$1 AND end_time IS NULL RETURNING id`
			_, err = tx.Exec(updateDefaultOfferingsStmt, workspaceId, now)
			if err != nil {
				log.Printf("PostgresDBStore.CreateAssignment: error updating default future offerings: %v\n", err)
				return "", err
			}
			log.Println("--- Cancelled all default Offerings ")

			// update any non-default offerings (cancel them)
			updateOfferingsStmt := `UPDATE offerings SET cancelled=TRUE WHERE workspace_id=$1 AND end_time >= $2 RETURNING id`
			_, err = tx.Exec(updateOfferingsStmt, workspaceId, now)
			if err != nil {
				log.Printf("PostgresDBStore.CreateDefaultOffering: error updating future offerings: %v\n", err)
				return "", err
			}
			log.Println("--- Cancelled all Offerings ")
			// create assignment
			var waId string
			createAssignmentStmt :=
				`INSERT INTO workspace_assignee(user_id, workspace_id, start_time) VALUES ($1, $2, $3) RETURNING id`
			err = tx.QueryRow(createAssignmentStmt, userId, workspaceId, now).Scan(&waId)
			if err != nil {
				return "", nil
			}
			log.Println("--- Created Assignment: ", waId)
		}
	}
	err = tx.Commit()
	return workspaceId, err
}

func (p PostgresDBStore) GetDeletedWorkspaces() ([]*model.Workspace, error) {
	rows, err := p.database.Query(`SELECT id, name, floor_id, details, metadata FROM workspaces WHERE deleted=TRUE;`)
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

func (p PostgresDBStore) DeleteWorkspaces(ids []string) error {
	tx, err := p.database.Begin()
	defer tx.Rollback()
	if err != nil {
		return err
	}
	for _, id := range ids {
		_, err := tx.Exec(`DELETE from workspaces where id=$1`, id)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
