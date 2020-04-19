package postgres

import (
	"errors"
	"go-api/model"
	"log"
	"time"
)

func (p PostgresDBStore) GetOneBooking(id string) (*model.Booking, error) {
	sqlStatement := `SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by FROM bookings WHERE id=$1;`
	var booking model.Booking
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&booking.ID,
		&booking.UserID,
		&booking.WorkspaceID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.Cancelled,
		&booking.CreatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (p PostgresDBStore) GetOneExpandedBooking(id string) (*model.ExpandedBooking, error) {
	sqlStatement := `SELECT b.id, u.id, w.id, b.start_time, b.end_time, b.cancelled, b.created_by, w.name, u.name, f.id, f.name
					 FROM bookings AS b
         			 INNER JOIN users AS u ON b.user_id = u.id
         			 INNER JOIN workspaces AS w ON b.workspace_id = w.id
         			 INNER JOIN floors AS f ON w.floor_id = f.id  
					 WHERE b.id=$1;`
	var eBooking model.ExpandedBooking
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&eBooking.ID,
		&eBooking.UserID,
		&eBooking.WorkspaceID,
		&eBooking.StartDate,
		&eBooking.EndDate,
		&eBooking.Cancelled,
		&eBooking.CreatedBy,
		&eBooking.WorkspaceName,
		&eBooking.UserName,
		&eBooking.FloorID,
		&eBooking.FloorName,
	)
	if err != nil {
		return nil, err
	}
	return &eBooking, nil
}

func (p PostgresDBStore) GetAllBookings() ([]*model.Booking, error) {
	sqlStatement := `SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by FROM bookings;`
	return p.queryMultipleBookings(sqlStatement)
}

func (p PostgresDBStore) GetAllExpandedBookings() ([]*model.ExpandedBooking, error) {
	sqlStatement := `SELECT b.id, u.id, w.id, b.start_time, b.end_time, b.cancelled, b.created_by, w.name, u.name, f.id, f.name
					 FROM bookings AS b
         			 INNER JOIN users AS u ON b.user_id = u.id
         			 INNER JOIN workspaces AS w ON b.workspace_id = w.id
         			 INNER JOIN floors AS f ON w.floor_id = f.id  
					 `
	return p.queryMultipleExpandedBookings(sqlStatement)
}

func (p PostgresDBStore) GetBookingsByWorkspaceID(id string) ([]*model.Booking, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by FROM bookings WHERE workspace_id=$1;`
	return p.queryMultipleBookings(sqlStatement, id)
}

func (p PostgresDBStore) GetExpandedBookingsByWorkspaceID(id string) ([]*model.ExpandedBooking, error) {
	sqlStatement :=
		`SELECT b.id, u.id, w.id, b.start_time, b.end_time, b.cancelled, b.created_by, w.name, u.name, f.id, f.name
		 FROM bookings AS b
		 INNER JOIN users AS u ON b.user_id = u.id
		 INNER JOIN workspaces AS w ON b.workspace_id = w.id
		 INNER JOIN floors AS f ON w.floor_id = f.id  
		 WHERE workspace_id=$1;`

	return p.queryMultipleExpandedBookings(sqlStatement, id)
}

func (p PostgresDBStore) GetBookingsByUserID(id string) ([]*model.Booking, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by FROM bookings WHERE user_id=$1;`
	return p.queryMultipleBookings(sqlStatement, id)
}

func (p PostgresDBStore) GetExpandedBookingsByUserID(id string) ([]*model.ExpandedBooking, error) {
	sqlStatement :=
		`SELECT b.id, u.id, w.id, b.start_time, b.end_time, b.cancelled, b.created_by, w.name, u.name, f.id, f.name
		 FROM bookings AS b
		 INNER JOIN users AS u ON b.user_id = u.id
		 INNER JOIN workspaces AS w ON b.workspace_id = w.id
		 INNER JOIN floors AS f ON w.floor_id = f.id  
		 WHERE u.id=$1;`

	return p.queryMultipleExpandedBookings(sqlStatement, id)
}

func (p PostgresDBStore) GetBookingsByDateRange(start time.Time, end time.Time) ([]*model.Booking, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by FROM bookings 
				WHERE (start_time >= $1 AND end_time <= $2) OR 
						(start_time <= $1 AND end_time >= $2) OR 
						(start_time <= $1 AND end_time >= $1) OR 
						(start_time <= $2 AND end_time >= $2);`
	return p.queryMultipleBookings(sqlStatement, start, end)
}

func (p PostgresDBStore) GetExpandedBookingsByDateRange(start time.Time, end time.Time) ([]*model.ExpandedBooking, error) {
	sqlStatement :=
		`SELECT b.id, u.id, w.id, b.start_time, b.end_time, b.cancelled, b.created_by, w.name, u.name, f.id, f.name
		 FROM bookings AS b
		 INNER JOIN users AS u ON b.user_id = u.id
		 INNER JOIN workspaces AS w ON b.workspace_id = w.id
		 INNER JOIN floors AS f ON w.floor_id = f.id  
		 WHERE (start_time >= $1 AND end_time <= $2) OR 
						(start_time <= $1 AND end_time >= $2) OR 
						(start_time <= $1 AND end_time >= $1) OR 
						(start_time <= $2 AND end_time >= $2);`

	return p.queryMultipleExpandedBookings(sqlStatement, start, end)
}

func (p PostgresDBStore) CreateBooking(booking *model.Booking) (string, error) {
	tx, err := p.database.Begin()
	if err != nil {
		return "", nil
	}
	// Check if offering still exists
	var count int
	err = tx.QueryRow(
		`SELECT count(*) FROM offerings 
					WHERE workspace_id=$1 AND cancelled=FALSE AND
                    	   (start_time <= $2 AND (end_time >= $3 OR end_time IS NULL))`,
		booking.WorkspaceID, booking.StartDate, booking.EndDate,
	).Scan(&count)
	if err != nil || count == 0 {
		return "", errors.New("invalid operation: workspace is not offered")
	}

	// Check if user has a booking at this time
	err = tx.QueryRow(
		`SELECT count(*) FROM bookings 
					WHERE user_id=$1 AND cancelled=FALSE AND
                    	   ((start_time <= $2 AND end_time >= $2) OR
                    	   (start_time <= $3 AND end_time >= $3) OR
                    	   (start_time >= $2 AND end_time <= $3) OR
                    	   (start_time <= $2 AND end_time >= $3))
                    	   `,
		booking.UserID, booking.StartDate, booking.EndDate,
	).Scan(&count) // if sql query fails, count won't be "updated" -> check error
	if err != nil || count != 0 {
		return "", errors.New("invalid operation: the user has a booking in this date range")
	}

	// Max 10 Bookings per user
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return "", errors.New("invalid operation: location failed")
	}
	timeNow := time.Now().In(loc)
	startOfDay := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, loc)
	err = tx.QueryRow(
		`SELECT count(*) FROM bookings 
					WHERE user_id=$1 AND
					cancelled=FALSE AND
					start_time >=$2`,
		booking.UserID, startOfDay,
	).Scan(&count)
	if err != nil || count > 10 {
		return "", errors.New("invalid operation: the user has 10 active bookings")
	}

	// Check for conflicts
	err = tx.QueryRow(
		`SELECT count(*) FROM bookings 
					WHERE workspace_id=$1 AND cancelled=FALSE AND
                    	   ((start_time <= $2 AND end_time >= $3) OR
                    	    (start_time <= $2 AND end_time >= $2) OR 
                    	    (start_time <= $3 AND end_time >= $3) OR
                    	    (start_time >= $2 AND end_time <= $3))`,
		booking.WorkspaceID, booking.StartDate, booking.EndDate,
	).Scan(&count)
	if err != nil || count > 0 {
		return "", errors.New("invalid operation: workspace already booked for this duration")
	}

	sqlStatement :=
		`INSERT INTO bookings(user_id, workspace_id, start_time, end_time, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id string
	err = tx.QueryRow(sqlStatement,
		booking.UserID,
		booking.WorkspaceID,
		booking.StartDate,
		booking.EndDate,
		booking.CreatedBy,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, tx.Commit()
}

func (p PostgresDBStore) UpdateBooking(id string, booking *model.Booking) error {
	sqlStatement :=
		`UPDATE bookings
				SET user_id = $2, workspace_id = $3, cancelled = $4, start_time = $5, end_time = $6, created_by = $7
				WHERE id = $1
				RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		id,
		booking.UserID,
		booking.WorkspaceID,
		booking.Cancelled,
		booking.StartDate,
		booking.EndDate,
		booking.CreatedBy,
	).Scan(&_id)
	if err != nil {
		return err
	}
	if _id != id {
		return CreateError
	}
	return nil
}

func (p PostgresDBStore) RemoveBooking(id string) error {
	sqlStatement :=
		`UPDATE bookings
				SET cancelled = true
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

func (p PostgresDBStore) GetExpiredBookings(since time.Time) ([]*model.Booking, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by FROM bookings 
				WHERE end_time < $1`
	return p.queryMultipleBookings(sqlStatement, since)
}

func (p PostgresDBStore) queryMultipleBookings(sqlStatement string, args ...interface{}) ([]*model.Booking, error) {
	rows, err := p.database.Query(sqlStatement, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bookings := make([]*model.Booking, 0)
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.WorkspaceID,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Cancelled,
			&booking.CreatedBy,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.queryMultipleBookings: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		bookings = append(bookings, &booking)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (p PostgresDBStore) queryMultipleExpandedBookings(sqlStatement string, args ...interface{}) ([]*model.ExpandedBooking, error) {
	rows, err := p.database.Query(sqlStatement, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bookings := make([]*model.ExpandedBooking, 0)
	for rows.Next() {
		var eBooking model.ExpandedBooking
		err := rows.Scan(
			&eBooking.ID,
			&eBooking.UserID,
			&eBooking.WorkspaceID,
			&eBooking.StartDate,
			&eBooking.EndDate,
			&eBooking.Cancelled,
			&eBooking.CreatedBy,
			&eBooking.WorkspaceName,
			&eBooking.UserName,
			&eBooking.FloorID,
			&eBooking.FloorName,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.queryMultipleBookings: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		bookings = append(bookings, &eBooking)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (p PostgresDBStore) DeleteBookings(ids []string) error {
	tx, err := p.database.Begin()
	defer tx.Rollback()
	if err != nil {
		return err
	}
	for _, id := range ids {
		_, err := tx.Exec(`DELETE from bookings where id=$1`, id)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
