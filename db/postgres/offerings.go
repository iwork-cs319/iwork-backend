package postgres

import (
	"go-api/model"
	"log"
	"time"
)

func (p PostgresDBStore) GetOneOffering(id string) (*model.Offering, error) {
	sqlStatement := `SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM offerings WHERE id=$1;`
	var offering model.Offering
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&offering.ID,
		&offering.UserID,
		&offering.WorkspaceID,
		&offering.StartDate,
		&offering.EndDate,
		&offering.Cancelled,
	)
	if err != nil {
		return nil, err
	}
	return &offering, nil
}

func (p PostgresDBStore) GetAllOfferings() ([]*model.Offering, error) {
	sqlStatement := `SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM offerings;`
	return p.queryMultipleOfferings(sqlStatement)
}

func (p PostgresDBStore) GetOfferingsByWorkspaceID(id string) ([]*model.Offering, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM offerings WHERE workspace_id=$1;`
	return p.queryMultipleOfferings(sqlStatement, id)
}

func (p PostgresDBStore) GetOfferingsByUserID(id string) ([]*model.Offering, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM offerings WHERE user_id=$1;`
	return p.queryMultipleOfferings(sqlStatement, id)
}

func (p PostgresDBStore) GetOfferingsByDateRange(start time.Time, end time.Time) ([]*model.Offering, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM offerings 
				WHERE start_time >= $1 AND end_time <= $2;`
	return p.queryMultipleOfferings(sqlStatement, start, end)
}

func (p PostgresDBStore) CreateOffering(offering *model.Offering) (string, error) {
	sqlStatement :=
		`INSERT INTO offerings(user_id, workspace_id, start_time, end_time) VALUES ($1, $2, $3, $4) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement,
		offering.UserID,
		offering.WorkspaceID,
		offering.StartDate,
		offering.EndDate,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p PostgresDBStore) UpdateOffering(id string, offering *model.Offering) error {
	sqlStatement :=
		`UPDATE offerings
				SET user_id = $2, workspace_id = $3, cancelled = $4, start_time = $5, end_time = $6
				WHERE id = $1
				RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		id,
		offering.UserID,
		offering.UserID,
		offering.Cancelled,
		offering.StartDate,
		offering.EndDate,
	).Scan(&_id)
	if err != nil {
		return err
	}
	if _id != id {
		return CreateError
	}
	return nil
}

func (p PostgresDBStore) RemoveOffering(id string) error {
	return nil
	//	sqlStatement :=
	//		`DELETE FROM offerings
	//				WHERE id = $1
	//				RETURNING id;`
	//	var _id string
	//	err := p.database.QueryRow(sqlStatement,
	//		id,
	//	).Scan(&_id)
	//	if err != nil {
	//		return err
	//	}
	//	if _id != id {
	//		return CreateError
	//	}
	//	return nil
}

func (p PostgresDBStore) queryMultipleOfferings(sqlStatement string, args ...interface{}) ([]*model.Offering, error) {
	rows, err := p.database.Query(sqlStatement, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	offerings := make([]*model.Offering, 0)
	for rows.Next() {
		var offering model.Offering
		err := rows.Scan(
			&offering.ID,
			&offering.UserID,
			&offering.WorkspaceID,
			&offering.StartDate,
			&offering.EndDate,
			&offering.Cancelled,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.queryMultipleOfferings: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		offerings = append(offerings, &offering)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return offerings, nil
}
