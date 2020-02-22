package postgres

import (
	"go-api/model"
	"log"
	"time"
)

func (p PostgresDBStore) GetOneBooking(id string) (*model.Booking, error) {
	sqlStatement := `SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM bookings WHERE id=$1;`
	var booking model.Booking
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&booking.ID,
		&booking.UserID,
		&booking.WorkspaceID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.Cancelled,
	)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (p PostgresDBStore) GetAllBookings() ([]*model.Booking, error) {
	sqlStatement := `SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM bookings;`
	return p.queryMultipleBookings(sqlStatement)
}

func (p PostgresDBStore) GetBookingsByWorkspaceID(id string) ([]*model.Booking, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM bookings WHERE workspace_id=$1;`
	return p.queryMultipleBookings(sqlStatement, id)
}

func (p PostgresDBStore) GetBookingsByUserID(id string) ([]*model.Booking, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM bookings WHERE user_id=$1;`
	return p.queryMultipleBookings(sqlStatement, id)
}

func (p PostgresDBStore) GetBookingsByDateRange(start time.Time, end time.Time) ([]*model.Booking, error) {
	sqlStatement :=
		`SELECT id, user_id, workspace_id, start_time, end_time, cancelled FROM bookings 
				WHERE start_time >= $1 AND end_time <= $2;`
	return p.queryMultipleBookings(sqlStatement, start, end)
}

func (p PostgresDBStore) CreateBooking(booking *model.Booking) (string, error) {
	sqlStatement :=
		`INSERT INTO bookings(user_id, workspace_id, start_time, end_time) VALUES ($1, $2, $3, $4) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement,
		booking.UserID,
		booking.WorkspaceID,
		booking.StartDate,
		booking.EndDate,
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
	panic("implement me")
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
