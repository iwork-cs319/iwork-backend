package postgres

import (
	"errors"
	"go-api/model"
	"go-api/utils"
	"log"
	"time"
)

func (p PostgresDBStore) GetOneOffering(id string) (*model.Offering, error) {
	sqlStatement := `SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by FROM offerings WHERE id=$1;`
	var offering model.Offering
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&offering.ID,
		&offering.UserID,
		&offering.WorkspaceID,
		&offering.StartDate,
		&offering.EndDate,
		&offering.Cancelled,
		&offering.CreatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &offering, nil
}

func (p PostgresDBStore) GetOneExpandedOffering(id string) (*model.ExpandedOffering, error) {
	sqlStatement := `SELECT o.id, u.id, w.id, o.start_time, o.end_time, o.cancelled, o.created_by, w.name, u.name, f.id, f.name
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
		&eOffering.CreatedBy,
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
	sqlStatement := `SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by FROM offerings;`
	return p.queryMultipleOfferings(sqlStatement)
}

func (p PostgresDBStore) GetAllExpandedOfferings() ([]*model.ExpandedOffering, error) {
	sqlStatement := `SELECT o.id, u.id, w.id, o.start_time, o.end_time, o.cancelled, o.created_by, w.name, u.name, f.id, f.name
					 FROM offerings AS o
         			 INNER JOIN users AS u ON o.user_id = u.id
         			 INNER JOIN workspaces AS w ON o.workspace_id = w.id
         			 INNER JOIN floors AS f ON w.floor_id = f.id  
					 `
	return p.queryMultipleExpandedOfferings(sqlStatement)
}

func (p PostgresDBStore) GetOfferingsByWorkspaceID(id string) ([]*model.Offering, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by FROM offerings WHERE workspace_id=$1;`
	return p.queryMultipleOfferings(sqlStatement, id)
}

func (p PostgresDBStore) GetExpandedOfferingsByWorkspaceID(id string) ([]*model.ExpandedOffering, error) {
	sqlStatement := `SELECT o.id, u.id, w.id, o.start_time, o.end_time, o.cancelled, o.created_by, w.name, u.name, f.id, f.name
					 FROM offerings AS o
         			 INNER JOIN users AS u ON o.user_id = u.id
         			 INNER JOIN workspaces AS w ON o.workspace_id = w.id
         			 INNER JOIN floors AS f ON w.floor_id = f.id
					 WHERE workspace_id=$1;
					 `
	return p.queryMultipleExpandedOfferings(sqlStatement, id)
}

func (p PostgresDBStore) GetOfferingsByUserID(id string) ([]*model.Offering, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by FROM offerings WHERE user_id=$1;`
	return p.queryMultipleOfferings(sqlStatement, id)
}

func (p PostgresDBStore) GetExpandedOfferingsByUserID(id string) ([]*model.ExpandedOffering, error) {
	sqlStatement := `SELECT o.id, u.id, w.id, o.start_time, o.end_time, o.cancelled, o.created_by, w.name, u.name, f.id, f.name
					 FROM offerings AS o
         			 INNER JOIN users AS u ON o.user_id = u.id
         			 INNER JOIN workspaces AS w ON o.workspace_id = w.id
         			 INNER JOIN floors AS f ON w.floor_id = f.id
					 WHERE u.id=$1;
					 `
	return p.queryMultipleExpandedOfferings(sqlStatement, id)
}

func (p PostgresDBStore) GetOfferingsByDateRange(start time.Time, end time.Time) ([]*model.Offering, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by FROM offerings 
				WHERE start_time <= $1 AND end_time >= $2;`
	return p.queryMultipleOfferings(sqlStatement, start, end)
}

func (p PostgresDBStore) GetExpandedOfferingsByDateRange(start time.Time, end time.Time) ([]*model.ExpandedOffering, error) {
	sqlStatement := `SELECT o.id, u.id, w.id, o.start_time, o.end_time, o.cancelled, o.created_by, w.name, u.name, f.id, f.name
					 FROM offerings AS o
         			 INNER JOIN users AS u ON o.user_id = u.id
         			 INNER JOIN workspaces AS w ON o.workspace_id = w.id
         			 INNER JOIN floors AS f ON w.floor_id = f.id
					 WHERE start_time <= $1 AND end_time >= $2;
					 `
	return p.queryMultipleExpandedOfferings(sqlStatement, start, end)
}

func (p PostgresDBStore) GetOfferingsByWorkspaceIDAndDateRange(id string, start time.Time, end time.Time) (*model.Offering, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by
		 FROM offerings
		 WHERE workspace_id=$1 AND start_time <= $2 AND end_time >= $3;`
	var offering model.Offering
	row := p.database.QueryRow(sqlStatement, id, start, end)
	err := row.Scan(
		&offering.ID,
		&offering.UserID,
		&offering.WorkspaceID,
		&offering.StartDate,
		&offering.EndDate,
		&offering.Cancelled,
		&offering.CreatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &offering, nil
}

func (p PostgresDBStore) CreateOffering(offering *model.Offering) (string, error) {
	tx, err := p.database.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Check if its currently assigned
	var count int
	err = tx.QueryRow(
		`SELECT count(*) FROM workspace_assignee 
					WHERE workspace_id=$1 AND 
                    	   (start_time <= $2 AND (end_time >= $3 OR end_time IS NULL))`,
		offering.WorkspaceID, offering.StartDate, offering.EndDate,
	).Scan(&count)
	if err != nil || count == 0 {
		return "", errors.New("invalid operation: unassigned workspace cannot be offered")
	}

	// Check for conflicts
	err = tx.QueryRow(
		`SELECT count(*) FROM offerings 
					WHERE workspace_id=$1 AND cancelled=FALSE AND
                    	   ((start_time <= $2 AND end_time >= $3) OR
                    	    (start_time <= $2 AND end_time >= $2) OR 
                    	    (start_time <= $3 AND end_time >= $3) OR
                    	    (start_time >= $2 AND end_time <= $3))`,
		offering.WorkspaceID, offering.StartDate, offering.EndDate,
	).Scan(&count)
	if err != nil || count > 0 {
		return "", errors.New("invalid operation: workspace already offered for this duration")
	}

	sqlStatement :=
		`INSERT INTO offerings(user_id, workspace_id, start_time, end_time, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id string
	err = tx.QueryRow(sqlStatement,
		offering.UserID,
		offering.WorkspaceID,
		offering.StartDate,
		offering.EndDate,
		offering.CreatedBy,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, tx.Commit()
}

func (p PostgresDBStore) UpdateOffering(id string, offering *model.Offering) error {
	sqlStatement :=
		`UPDATE offerings
				SET user_id = $2, workspace_id = $3, cancelled = $4, start_time = $5, end_time = $6, created_by = $7
				WHERE id = $1
				RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		id,
		offering.UserID,
		offering.WorkspaceID,
		offering.Cancelled,
		offering.StartDate,
		offering.EndDate,
		offering.CreatedBy,
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
	tx, err := p.database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var userId, workspaceId string
	var start, end time.Time
	err = tx.QueryRow(`SELECT user_id, workspace_id, start_time, end_time FROM offerings WHERE id=$1`,
		id,
	).Scan(&userId, &workspaceId, &start, &end)
	if err != nil {
		return err
	}
	if userId == utils.EmptyUserUUID {
		return errors.New("invalid operation: cannot remove offerings by default user")
	}
	var count int
	err = tx.QueryRow(
		`SELECT count(*) FROM bookings WHERE workspace_id=$1 AND (start_time >= $2 AND end_time <= $3)`,
		workspaceId, start, end,
	).Scan(&count)
	if count > 0 {
		return errors.New("invalid operation: conflicting bookings for this offering period; cannot delete")
	}
	sqlStatement :=
		`UPDATE offerings
				SET cancelled = true
				WHERE id = $1
				RETURNING id;`
	var _id string
	err = tx.QueryRow(sqlStatement,
		id,
	).Scan(&_id)
	if err != nil {
		return err
	}
	if _id != id {
		return CreateError
	}
	return tx.Commit()
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
			&offering.CreatedBy,
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
	offerings := make([]*model.ExpandedOffering, 0)
	for rows.Next() {
		var eOffering model.ExpandedOffering
		err := rows.Scan(
			&eOffering.ID,
			&eOffering.UserID,
			&eOffering.WorkspaceID,
			&eOffering.StartDate,
			&eOffering.EndDate,
			&eOffering.Cancelled,
			&eOffering.CreatedBy,
			&eOffering.WorkspaceName,
			&eOffering.UserName,
			&eOffering.FloorID,
			&eOffering.FloorName,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.queryMultipleBookings: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		offerings = append(offerings, &eOffering)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return offerings, nil
}
