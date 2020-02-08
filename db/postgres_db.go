package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"go-api/model"
	"time"
)

type PostgresDBStore struct {
	database *sql.DB
}

func (p PostgresDBStore) GetOneWorkspace(id string) (*model.Workspace, error) {
	sqlStatement := `SELECT id, name, floor FROM users WHERE id=$1;`
	var workspace model.Workspace
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(&workspace.ID, &workspace.Name, &workspace.Floor)
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (p PostgresDBStore) UpdateWorkspace(id string, workspace *model.Workspace) error {
	panic("implement me")
}

func (p PostgresDBStore) CreateWorkspace(workspace *model.Workspace) error {
	panic("implement me")
}

func (p PostgresDBStore) RemoveWorkspace(id string) error {
	panic("implement me")
}

func (p PostgresDBStore) GetAllWorkspaces() ([]*model.Workspace, error) {
	panic("implement me")
}

func (p PostgresDBStore) GetOneBooking(id string) (*model.Booking, error) {
	panic("implement me")
}

func (p PostgresDBStore) GetAllBookings() ([]*model.Booking, error) {
	panic("implement me")
}

func (p PostgresDBStore) GetBookingsByWorkspaceID(id string) ([]*model.Booking, error) {
	panic("implement me")
}

func (p PostgresDBStore) GetBookingsByUserID(id string) ([]*model.Booking, error) {
	panic("implement me")
}

func (p PostgresDBStore) GetBookingsByDateRange(start time.Time, end time.Time) ([]*model.Booking, error) {
	panic("implement me")
}

func (p PostgresDBStore) CreateBooking(booking *model.Booking) error {
	panic("implement me")
}

func (p PostgresDBStore) UpdateBooking(id string, booking *model.Booking) error {
	panic("implement me")
}

func (p PostgresDBStore) RemoveBooking(id string) error {
	panic("implement me")
}

func (p PostgresDBStore) Close() {
	p.database.Close()
}

func NewPostgresDataStore(dbUrl string) (*DataStore, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DataStore{
		WorkspaceProvider: &PostgresDBStore{db},
		BookingProvider:   nil,
	}, nil
}
