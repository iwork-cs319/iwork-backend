package db

import (
	"errors"
	"time"
)
import "go-api/model"

var NotFoundError = errors.New("not found")
var EmptyError = errors.New("empty")

type DataStore struct {
	WorkspaceProvider workspaceProvider
	BookingProvider   bookingProvider
	OfferingProvider  offeringProvider
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
	GetAllBookings() ([]*model.Booking, error)
	GetBookingsByWorkspaceID(id string) ([]*model.Booking, error)
	GetBookingsByUserID(id string) ([]*model.Booking, error)
	GetBookingsByDateRange(start time.Time, end time.Time) ([]*model.Booking, error)
	CreateBooking(booking *model.Booking) error
	UpdateBooking(id string, booking *model.Booking) error
	RemoveBooking(id string) error
}

type offeringProvider interface {
	GetOneOffering(id string) (*model.Offering, error)
	GetAllOfferings() ([]*model.Offering, error)
	GetOfferingsByWorkspaceID(id string) ([]*model.Offering, error)
	GetOfferingsByUserID(id string) ([]*model.Offering, error)
	GetOfferingsByDateRange(start time.Time, end time.Time) ([]*model.Offering, error)
	CreateOffering(booking *model.Offering) error
	UpdateOffering(id string, booking *model.Offering) error
	RemoveOffering(id string) error
}
