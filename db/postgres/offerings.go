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

func (p PostgresDBStore) GetExpandedOneOffering(id string) (*model.ExpandedOffering, error) {
	sqlStatement := `SELECT o.id, u.id, w.id, o.start_time, o.end_time, o.cancelled, w.name, u.name, f.id, f.name
					 FROM offerings AS o
         			 INNER JOIN users AS u ON o.user_id = u.id
         			 INNER JOIN workspaces AS w ON o.workspace_id = w.id
         			 INNER JOIN floors AS f ON w.floor_id = f.id  
					 WHERE o.id=$1;`
	var eOffering model.ExpandedOffering
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&eOffering.ID,
		&eOffering.UserID,
		&eOffering.WorkspaceID,
		&eOffering.StartDate,
		&eOffering.EndDate,
		&eOffering.Cancelled,
		&eOffering.WorkspaceName,
		&eOffering.UserName,
		&eOffering.FloorID,
		&eOffering.FloorName,
	)
	if err != nil {
		return nil, err
	}
	return &eOffering, nil
}

func (p PostgresDBStore) GetAllOfferings() ([]*model.Offering, error) {
	sqlStatement := `SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM offerings;`
	return p.queryMultipleOfferings(sqlStatement)
}

func (p PostgresDBStore) GetAllExpandedOfferings() ([]*model.ExpandedOffering, error) {
	sqlStatement := `SELECT o.id, u.id, w.id, o.start_time, o.end_time, b.cancelled, w.name, u.name, f.id, f.name
					 FROM offerings AS o
         			 INNER JOIN users AS u ON b.user_id = u.id
         			 INNER JOIN workspaces AS w ON o.workspace_id = w.id
         			 INNER JOIN floors AS f ON w.floor_id = f.id  
					 `
	return p.queryMultipleExpandedOfferings(sqlStatement)
}

func (p PostgresDBStore) GetOfferingsByWorkspaceID(id string) ([]*model.Offering, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM offerings WHERE workspace_id=$1;`
	return p.queryMultipleOfferings(sqlStatement, id)
}

func (p PostgresDBStore) GetExpandedOfferingsByWorkspaceID() ([]*model.ExpandedOffering, error) {
	sqlStatement := `SELECT o.id, u.id, w.id, o.start_time, o.end_time, b.cancelled, w.name, u.name, f.id, f.name
					 FROM offerings AS o
         			 INNER JOIN users AS u ON b.user_id = u.id
         			 INNER JOIN workspaces AS w ON o.workspace_id = w.id
         			 INNER JOIN floors AS f ON w.floor_id = f.id
					 WHERE workspace_id=$1;
					 `
	return p.queryMultipleExpandedOfferings(sqlStatement)
}

func (p PostgresDBStore) GetOfferingsByUserID(id string) ([]*model.Offering, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM offerings WHERE user_id=$1;`
	return p.queryMultipleOfferings(sqlStatement, id)
}

func (p PostgresDBStore) GetExpandedOfferingsByUserID() ([]*model.ExpandedOffering, error) {
	sqlStatement := `SELECT o.id, u.id, w.id, o.start_time, o.end_time, b.cancelled, w.name, u.name, f.id, f.name
					 FROM offerings AS o
         			 INNER JOIN users AS u ON b.user_id = u.id
         			 INNER JOIN workspaces AS w ON o.workspace_id = w.id
         			 INNER JOIN floors AS f ON w.floor_id = f.id
					 WHERE u.id=$1;
					 `
	return p.queryMultipleExpandedOfferings(sqlStatement)
}

func (p PostgresDBStore) GetOfferingsByDateRange(start time.Time, end time.Time) ([]*model.Offering, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM offerings 
				WHERE start_time >= $1 AND end_time <= $2;`
	return p.queryMultipleOfferings(sqlStatement, start, end)
}

func (p PostgresDBStore) GetExpandedOfferingsByDateRange() ([]*model.ExpandedOffering, error) {
	sqlStatement := `SELECT o.id, u.id, w.id, o.start_time, o.end_time, b.cancelled, w.name, u.name, f.id, f.name
					 FROM offerings AS o
         			 INNER JOIN users AS u ON b.user_id = u.id
         			 INNER JOIN workspaces AS w ON o.workspace_id = w.id
         			 INNER JOIN floors AS f ON w.floor_id = f.id
					 WHERE start_time >= $1 AND end_time <= $2;
					 `
	return p.queryMultipleExpandedOfferings(sqlStatement)
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
	sqlStatement :=
		`UPDATE offerings
				SET cancelled = false
				WHERE id = $1
				RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		id,
	).Scan(&_id)
	if err != nil {
		return err
	}
	if _id != id {
		return CreateError
	}
	return nil
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

func (p PostgresDBStore) queryMultipleExpandedOfferings(sqlStatement string, args ...interface{}) ([]*model.ExpandedOffering, error) {
	rows, err := p.database.Query(sqlStatement, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bookings := make([]*model.ExpandedOffering, 0)
	for rows.Next() {
		var eOffering model.ExpandedOffering
		err := rows.Scan(
			&eOffering.ID,
			&eOffering.UserID,
			&eOffering.WorkspaceID,
			&eOffering.StartDate,
			&eOffering.EndDate,
			&eOffering.Cancelled,
			&eOffering.WorkspaceName,
			&eOffering.UserName,
			&eOffering.FloorID,
			&eOffering.FloorName,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.queryMultipleBookings: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		bookings = append(bookings, &eOffering)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return bookings, nil
}
