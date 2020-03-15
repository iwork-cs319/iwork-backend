package db

import (
	"errors"
	"go-api/model"
	"time"
)

var NotFoundError = errors.New("not found")
var EmptyError = errors.New("empty")

type DataStore struct {
	Closable
	WorkspaceProvider workspaceProvider
	BookingProvider   bookingProvider
	UserProvider      userProvider
	FloorProvider     floorProvider
	OfferingProvider  offeringProvider
}

type Closable interface {
	Close()
}

type workspaceProvider interface {
	GetOneWorkspace(id string) (*model.Workspace, error)
	UpdateWorkspace(id string, workspace *model.Workspace) error
	CreateWorkspace(workspace *model.Workspace) (string, error)
	RemoveWorkspace(id string) error
	GetAllWorkspaces() ([]*model.Workspace, error)
	GetAllWorkspacesByFloor(floorId string) ([]*model.Workspace, error)
	FindAvailability(floorId string, start time.Time, end time.Time) ([]string, error)
	CreateAssignment(userId, workspaceId string) error
}

type bookingProvider interface {
	GetOneBooking(id string) (*model.Booking, error)
	GetOneExpandedBooking(id string) (*model.ExpandedBooking, error)
	GetAllBookings() ([]*model.Booking, error)
	GetAllExpandedBookings() ([]*model.ExpandedBooking, error)
	GetBookingsByWorkspaceID(id string) ([]*model.Booking, error)
	GetExpandedBookingsByWorkspaceID(id string) ([]*model.ExpandedBooking, error)
	GetBookingsByUserID(id string) ([]*model.Booking, error)
	GetExpandedBookingsByUserID(id string) ([]*model.ExpandedBooking, error)
	GetBookingsByDateRange(start time.Time, end time.Time) ([]*model.Booking, error)
	GetExpandedBookingsByDateRange(start time.Time, end time.Time) ([]*model.ExpandedBooking, error)
	CreateBooking(booking *model.Booking) (string, error)
	UpdateBooking(id string, booking *model.Booking) error
	RemoveBooking(id string) error
}

type userProvider interface {
	GetOneUser(id string) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	CreateUser(user *model.User) error
	//UpdateUser(id string, user *model.User) error
	//RemoveUser(id string) error
}

type floorProvider interface {
	GetOneFloor(id string) (*model.Floor, error)
	GetAllFloors() ([]*model.Floor, error)
	CreateFloor(floor *model.Floor) (string, error)
	//UpdateFloor(id string, user *model.Floor) error
	//RemoveFloor(id string) error
}

type offeringProvider interface {
	GetOneOffering(id string) (*model.Offering, error)
	GetOneExpandedOffering(id string) (*model.ExpandedOffering, error)
	GetAllOfferings() ([]*model.Offering, error)
	GetAllExpandedOfferings() ([]*model.ExpandedOffering, error)
	GetOfferingsByWorkspaceID(id string) ([]*model.Offering, error)
	GetExpandedOfferingsByWorkspaceID(id string) ([]*model.ExpandedOffering, error)
	GetOfferingsByUserID(id string) ([]*model.Offering, error)
	GetExpandedOfferingsByUserID(id string) ([]*model.ExpandedOffering, error)
	GetOfferingsByDateRange(start time.Time, end time.Time) ([]*model.Offering, error)
	GetExpandedOfferingsByDateRange(start time.Time, end time.Time) ([]*model.ExpandedOffering, error)
	CreateOffering(booking *model.Offering) (string, error)
	UpdateOffering(id string, booking *model.Offering) error
	RemoveOffering(id string) error
}
