package db

import "errors"
import "go-api/model"

var NotFoundError = errors.New("not found")
var EmptyError = errors.New("empty")

type DataStore struct {
	WorkspaceProvider workspaceProvider
	//BookingProvider   bookingProvider
}

type workspaceProvider interface {
	GetOneWorkspace(id string) (*model.Workspace, error)
	UpdateWorkspace(id string, workspace *model.Workspace) error
	CreateWorkspace(workspace *model.Workspace) error
	RemoveWorkspace(id string) error
	GetAllWorkspaces() ([]*model.Workspace, error)
}

type bookingProvider interface {
	GetOneBooking(id string) (*model.Booking, error)
	GetAllBookings(id string) (*model.Booking, error)
	GetBookingByWorkspaceID(id string) (*[]model.Booking, error)
	GetBookingByUserID(id string) (*[]model.Booking, error)
	CreateBooking(booking *model.Booking) error
	UpdateBooking(id string, booking *model.Booking) error
}
