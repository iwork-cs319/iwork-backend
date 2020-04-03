package postgres

import (
	"errors"
	"go-api/model"
	"log"
)

func (p PostgresDBStore) GetOneFloor(id string) (*model.Floor, error) {
	sqlStatement := `SELECT id, name, download_url, address FROM floors WHERE id=$1;`
	var floor model.Floor
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&floor.ID,
		&floor.Name,
		&floor.DownloadURL,
		&floor.Address,
	)
	if err != nil {
		return nil, err
	}
	return &floor, nil
}

func (p PostgresDBStore) GetAllFloors() ([]*model.Floor, error) {
	sqlStatement := `SELECT id, name, download_url, address FROM floors WHERE deleted=FALSE;`
	return p.queryMultipleFloors(sqlStatement)
}

func (p PostgresDBStore) GetAllFloorIDs() ([]string, error) {
	sqlStatement := `SELECT id FROM floors WHERE deleted=FALSE;`
	rows, err := p.database.Query(sqlStatement)
	if err != nil {
		log.Printf("PostgresDBStore.GetAllFloorIDs: %v, sqlStatement: %s\n", err, sqlStatement)
		return nil, err
	}
	defer rows.Close()
	floorIDs := make([]string, 0)
	for rows.Next() {
		var id string
		err := rows.Scan(
			&id,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.GetAllFloorIDs: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		floorIDs = append(floorIDs, id)
	}
	return floorIDs, nil
}

func (p PostgresDBStore) CreateFloor(floor *model.Floor) (string, error) {
	sqlStatement :=
		`INSERT INTO floors(name, download_url, address) VALUES ($1, $2, $3) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement,
		floor.Name,
		floor.DownloadURL,
		floor.Address,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

//func (p PostgresDBStore) UpdateFloor(id string, floor *model.Floor) error {
//	sqlStatement :=
//		`UPDATE floors
//				SET name = $2,
//				WHERE id = $1
//				RETURNING id;`
//	var _id string
//	err := p.database.QueryRow(sqlStatement,
//		id,
//		floor.Name,
//	).Scan(&_id)
//	if err != nil {
//		return err
//	}
//	if _id != id {
//		return CreateError
//	}
//	return nil
//}
//

func (p PostgresDBStore) RemoveFloor(id string, force bool) error {
	tx, err := p.database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	existingBookings := `SELECT count(*) FROM bookings b INNER JOIN workspaces w on b.workspace_id = w.id
							WHERE w.floor_id=$1 AND (b.start_time >= $2 OR b.end_time >=$2)`
	var count int
	err = tx.QueryRow(existingBookings,
		id,
	).Scan(&count)

	if count > 0 && !force {
		log.Println("Postgres.RemoveFloor: there are existing bookings for workspaces on this floor")
		return errors.New("invalid operation: there are existing bookings for workspaces on this floor")
	}

	deleteWorkspaceStmt := `UPDATE workspaces SET deleted=true WHERE floor_id=$1`
	_, err = tx.Exec(deleteWorkspaceStmt, id)
	if err != nil {
		log.Println("Postgres.RemoveFloor: error setting deleted on workspace")
		return err
	}
	deletedStmt :=
		`UPDATE floors SET deleted=true WHERE id = $1`
	_, err = tx.Exec(deletedStmt, id)
	if err != nil {
		log.Println("Postgres.RemoveFloor: error setting deleted on floor")
		return err
	}
	return tx.Commit()
}

func (p PostgresDBStore) queryMultipleFloors(sqlStatement string, args ...interface{}) ([]*model.Floor, error) {
	rows, err := p.database.Query(sqlStatement, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	floors := make([]*model.Floor, 0)
	for rows.Next() {
		var floor model.Floor
		err := rows.Scan(
			&floor.ID,
			&floor.Name,
			&floor.DownloadURL,
			&floor.Address,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.queryMultipleFloors: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		floors = append(floors, &floor)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return floors, nil
}
