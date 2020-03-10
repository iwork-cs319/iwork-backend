package postgres

import (
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
				WHERE start_time >= $1 AND end_time <= $2;`
	return p.queryMultipleBookings(sqlStatement, start, end)
}

func (p PostgresDBStore) GetExpandedBookingsByDateRange(start time.Time, end time.Time) ([]*model.ExpandedBooking, error) {
	sqlStatement :=
		`SELECT b.id, u.id, w.id, b.start_time, b.end_time, b.cancelled, b.created_by, w.name, u.name, f.id, f.name
		 FROM bookings AS b
		 INNER JOIN users AS u ON b.user_id = u.id
		 INNER JOIN workspaces AS w ON b.workspace_id = w.id
		 INNER JOIN floors AS f ON w.floor_id = f.id  
		 WHERE start_time >= $1 AND end_time <= $2;`

	return p.queryMultipleExpandedBookings(sqlStatement, start, end)
}

func (p PostgresDBStore) CreateBooking(booking *model.Booking) (string, error) {
	sqlStatement :=
		`INSERT INTO bookings(user_id, workspace_id, start_time, end_time, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement,
		booking.UserID,
		booking.WorkspaceID,
		booking.StartDate,
		booking.EndDate,
		booking.CreatedBy,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p PostgresDBStore) UpdateBooking(id string, booking *model.Booking) error {
	sqlStatement :=
		`UPDATE bookings
				SET user_id = $2, workspace_id = $3, cancelled = $4, start_time = $5, end_time = $6
				WHERE id = $1
				RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		id,
		booking.UserID,
		booking.UserID,
		booking.Cancelled,
		booking.StartDate,
		booking.EndDate,
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
			&eBooking.WorkspaceName,
			&eBooking.UserName,
			&eBooking.FloorID,
			&eBooking.FloorName,
			&eBooking.CreatedBy,
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
