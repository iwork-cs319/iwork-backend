package db

import (
	"go-api/model"
	"time"
)

type LocalDBStore struct {
	workspaces map[string]*model.Workspace
	bookings   map[string]*model.Booking
}

func (l LocalDBStore) GetOneWorkspace(id string) (*model.Workspace, error) {
	w, ok := l.workspaces[id]
	if !ok {
		return nil, NotFoundError
	}
	return w, nil
}

func (l LocalDBStore) UpdateWorkspace(id string, workspace *model.Workspace) error {
	_, ok := l.workspaces[id]
	if !ok {
		return NotFoundError
	}
	if workspace.Floor != "" {
		l.workspaces[id].Floor = workspace.Floor
	}
	if workspace.Name != "" {
		l.workspaces[id].Name = workspace.Name
	}
	if workspace.Props != nil {
		l.workspaces[id].Props = workspace.Props
	}
	return nil
}

func (l LocalDBStore) CreateWorkspace(workspace *model.Workspace) error {
	l.workspaces[workspace.ID] = workspace
	return nil
}

func (l LocalDBStore) GetAllWorkspaces() ([]*model.Workspace, error) {
	var list []*model.Workspace
	if len(l.workspaces) < 1 {
		return nil, EmptyError
	}
	for _, w := range l.workspaces {
		list = append(list, w)
	}
	return list, nil
}

func (l LocalDBStore) RemoveWorkspace(id string) error {
	delete(l.workspaces, id)
	return nil
}

func (l LocalDBStore) GetOneBooking(id string) (*model.Booking, error) {
	b, ok := l.bookings[id]
	if !ok {
		return nil, NotFoundError
	}
	return b, nil
}

func (l LocalDBStore) GetAllBookings() ([]*model.Booking, error) {
	var list []*model.Booking
	if len(l.bookings) < 1 {
		return nil, EmptyError
	}
	for _, w := range l.bookings {
		list = append(list, w)
	}
	return list, nil
}

func (l LocalDBStore) GetBookingsByWorkspaceID(id string) ([]*model.Booking, error) {
	var list []*model.Booking
	if len(l.bookings) < 1 {
		return nil, EmptyError
	}
	for _, b := range l.bookings {
		if b.WorkspaceID == id {
			list = append(list, b)
		}
	}
	return list, nil
}

func (l LocalDBStore) GetBookingsByUserID(id string) ([]*model.Booking, error) {
	var list []*model.Booking
	if len(l.bookings) < 1 {
		return nil, EmptyError
	}
	for _, b := range l.bookings {
		if b.UserID == id {
			list = append(list, b)
		}
	}
	return list, nil
}

func (l LocalDBStore) GetBookingsByDateRange(start time.Time, end time.Time) ([]*model.Booking, error) {
	var list []*model.Booking
	if len(l.bookings) < 1 {
		return nil, EmptyError
	}
	for _, b := range l.bookings {
		if b.StartDate.After(start) && b.EndDate.Before(end) { // Todo: Right way around?
			list = append(list, b)
		}
	}
	return list, nil
}

func (l LocalDBStore) CreateBooking(booking *model.Booking) error {
	l.bookings[booking.ID] = booking
	return nil
}

func (l LocalDBStore) UpdateBooking(id string, booking *model.Booking) error {
	_, ok := l.bookings[id]
	if !ok {
		return NotFoundError
	}
	if booking.WorkspaceID != "" {
		l.bookings[id].WorkspaceID = booking.WorkspaceID
	}
	if booking.UserID != "" {
		l.bookings[id].UserID = booking.WorkspaceID
	}
	if !booking.StartDate.IsZero() {
		l.bookings[id].StartDate = booking.StartDate
	}
	if !booking.EndDate.IsZero() {
		l.bookings[id].EndDate = booking.EndDate
	}
	//if booking.Canceled != nil {
	//	l.bookings[id].Canceled = booking.Canceled
	//}
	return nil
}

func (l LocalDBStore) RemoveBooking(id string) error {
	delete(l.bookings, id)
	return nil
}

func NewLocalDataStore() *DataStore {
	return &DataStore{
		WorkspaceProvider: &LocalDBStore{workspaces: map[string]*model.Workspace{
			"1": {
				ID:    "1",
				Name:  "Workspace 1",
				Props: nil,
			},
			"2": {
				ID:    "2",
				Name:  "Workspace 2",
				Props: nil,
			},
			"3": {
				ID:    "3",
				Name:  "Workspace 3",
				Props: nil,
			},
			"6": {
				ID:    "6",
				Name:  "Workspace 6",
				Props: nil,
			},
		}},
		BookingProvider: &LocalDBStore{bookings: map[string]*model.Booking{
			"1": {
				ID:          "1",
				WorkspaceID: "1",
				UserID:      "1",
				StartDate:   time.Unix(1580869576, 0),
				EndDate:     time.Unix(1580947199, 0),
				Cancelled:   false,
			},
			"2": {
				ID:          "2",
				WorkspaceID: "2",
				UserID:      "2",
				StartDate:   time.Unix(1571011200, 0),
				EndDate:     time.Unix(1571183999, 0),
				Cancelled:   true,
			},
		}},
	}
}
