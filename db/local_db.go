package db

import "go-api/model"

type LocalDBStore struct {
	workspaces map[string]*model.Workspace
	bookings map[string]*model.Booking
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

func (l LocalDBStore) RemoveWorkspace(id string) error {
	delete(l.workspaces, id)
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

func (l LocalDBStore) GetOneBooking(id string) (*model.Booking, error) {
	b, ok := l.bookings[id]
	if !ok {
		return nil, NotFoundError
	}
	return b, nil
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
	if booking.StartDate != nil {
		l.bookings[id].StartDate = booking.StartDate
	}
	if booking.EndDate != nil {
		l.bookings[id].EndDate = booking.EndDate
	}
	if booking.Canceled != nil {
		l.bookings[id].Canceled = booking.Canceled
	}
	return nil
}

func (l LocalDBStore) CreateBooking(booking *model.Booking) error {
	l.bookings[booking.ID] = booking
	return nil
}

func (l LocalDBStore) RemoveBooking(id string) error {
	delete(l.bookings, id)
	return nil
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
				ID:    "1",
				WorkspaceID:  "1",
				UserID: nil, // Not Implemented!
				StartDate: "2020/01/30",
				EndDate: "2020/01/30", // One day booking
				Canceled: false,
			},
			"2": {
				ID:    "2",
				WorkspaceID:  "2",
				UserID: nil, // Not Implemented!
				StartDate: "2020/01/28",
				EndDate: "2020/01/30", // Range booking
				Canceled: true, // Canceled
			},
		}},
	}
}
