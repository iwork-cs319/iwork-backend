package db

import (
	"go-api/model"
	"time"
)

type LocalDBStore struct {
	workspaces map[string]*model.Workspace
	bookings   map[string]*model.Booking
	offerings  map[string]*model.Offering
	users  map[string]*model.User
	floors  map[string]*model.Floor
}

func (l LocalDBStore) GetOneOffering(id string) (*model.Offering, error) {
	b, ok := l.offerings[id]
	if !ok {
		return nil, NotFoundError
	}
	return b, nil
}

func (l LocalDBStore) GetAllOfferings() ([]*model.Offering, error) {
	list := make([]*model.Offering, 0)
	if len(l.offerings) < 1 {
		return nil, EmptyError
	}
	for _, w := range l.offerings {
		list = append(list, w)
	}
	return list, nil
}

func (l LocalDBStore) GetOfferingsByWorkspaceID(id string) ([]*model.Offering, error) {
	list := make([]*model.Offering, 0)
	if len(l.offerings) < 1 {
		return nil, EmptyError
	}
	for _, b := range l.offerings {
		if b.WorkspaceID == id {
			list = append(list, b)
		}
	}
	return list, nil
}

func (l LocalDBStore) GetOfferingsByUserID(id string) ([]*model.Offering, error) {
	list := make([]*model.Offering, 0)
	if len(l.offerings) < 1 {
		return nil, EmptyError
	}
	for _, b := range l.offerings {
		if b.UserID == id {
			list = append(list, b)
		}
	}
	return list, nil
}

func (l LocalDBStore) GetOfferingsByDateRange(start time.Time, end time.Time) ([]*model.Offering, error) {
	list := make([]*model.Offering, 0)
	if len(l.offerings) < 1 {
		return nil, EmptyError
	}
	for _, b := range l.offerings {
		if b.StartDate.After(start) && b.EndDate.Before(end) { // Todo: Right way around?
			list = append(list, b)
		}
	}
	return list, nil
}

func (l LocalDBStore) CreateOffering(offering *model.Offering) error {
	l.offerings[offering.ID] = offering
	return nil
}

func (l LocalDBStore) UpdateOffering(id string, offering *model.Offering) error {
	_, ok := l.offerings[id]
	if !ok {
		return NotFoundError
	}
	if offering.WorkspaceID != "" {
		l.offerings[id].WorkspaceID = offering.WorkspaceID
	}
	if offering.UserID != "" {
		l.offerings[id].UserID = offering.WorkspaceID
	}
	if !offering.StartDate.IsZero() {
		l.offerings[id].StartDate = offering.StartDate
	}
	if !offering.EndDate.IsZero() {
		l.offerings[id].EndDate = offering.EndDate
	}
	return nil
}

func (l LocalDBStore) RemoveOffering(id string) error {
	delete(l.offerings, id)
	return nil
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
	list := make([]*model.Workspace, 0)
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
	list := make([]*model.Booking, 0)
	if len(l.bookings) < 1 {
		return nil, EmptyError
	}
	for _, w := range l.bookings {
		list = append(list, w)
	}
	return list, nil
}

func (l LocalDBStore) GetBookingsByWorkspaceID(id string) ([]*model.Booking, error) {
	list := make([]*model.Booking, 0)
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
	list := make([]*model.Booking, 0)
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
	list := make([]*model.Booking, 0)
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

func (l LocalDBStore) GetOneUser(id string) (*model.User, error) {
	u, ok := l.users[id]
	if !ok {
		return nil, NotFoundError
	}
	return u, nil
}

func (l LocalDBStore) GetAllUsers() ([]*model.User, error) {
	list := make([]*model.User, 0)
	if len(l.users) < 1 {
		return nil, EmptyError
	}
	for _, w := range l.users {
		list = append(list, w)
	}
	return list, nil
}

//func (l LocalDBStore) CreateUser(user *model.User) error {
//	l.users[user.ID] = user
//	return nil
//}

func (l LocalDBStore) UpdateUser(id string, user *model.User) error {
	_, ok := l.users[id]
	if !ok {
		return NotFoundError
	}
	//if user.WorkspaceID != "" {
	//	l.users[id].WorkspaceID = offering.WorkspaceID
	//}
	if user.Name != "" {
		l.users[id].Name = user.Name
	}
	if user.Department != "" {
		l.users[id].Department = user.Department
	}
	//if user.IsAdmin != nil {
	//	l.users[id].IsAdmin = user.IsAdmin
	//}
	return nil
}

//func (l LocalDBStore) RemoveUser(id string) error {
//	delete(l.users, id)
//	return nil
//}

func (l LocalDBStore) Close() {}

//func NewLocalDataStore() *DataStore {
//	localStore := &LocalDBStore{
//		bookings: map[string]*model.Booking{
//			"1": {
//				ID:          "1",
//				WorkspaceID: "1",
//				UserID:      "1",
//				StartDate:   time.Unix(1580869576, 0),
//				EndDate:     time.Unix(1580947199, 0),
//				Cancelled:   false,
//			},
//			"2": {
//				ID:          "2",
//				WorkspaceID: "2",
//				UserID:      "2",
//				StartDate:   time.Unix(1571011200, 0),
//				EndDate:     time.Unix(1571183999, 0),
//				Cancelled:   true,
//			},
//		},
//		workspaces: map[string]*model.Workspace{
//			"1": {
//				ID:    "1",
//				Name:  "Workspace 1",
//				Props: nil,
//			},
//			"2": {
//				ID:    "2",
//				Name:  "Workspace 2",
//				Props: nil,
//			},
//			"3": {
//				ID:    "3",
//				Name:  "Workspace 3",
//				Props: nil,
//			},
//			"6": {
//				ID:    "6",
//				Name:  "Workspace 6",
//				Props: nil,
//			},
//		},
//		offerings: map[string]*model.Offering{
//			"1": {
//				ID:          "1",
//				WorkspaceID: "1",
//				UserID:      "1",
//				StartDate:   time.Unix(1580859576, 0),
//				EndDate:     time.Unix(1580957199, 0),
//				Cancelled:   false,
//			},
//		},
//		users: map[string]*model.User{
//			"1": {
//				ID:          "1",
//				Name: "Prayansh",
//				Department: "IT",
//				IsAdmin: false,
//			},
//			"2": {
//				ID:          "2",
//				Name: "Ming",
//				Department: "IT",
//				IsAdmin: false,
//			},
//			"3": {
//				ID:          "3",
//				Name: "Alex",
//				Department: "IT",
//				IsAdmin: true,
//			},
//
//		},
//		floors: map[string]*model.Floor{
//			"1": {
//				ID:          "1",
//				Name: "Floor 1",
//			},
//			"2": {
//				ID:          "2",
//				Name: "Floor 99",
//			},
//
//		}}
//	return &DataStore{
//		WorkspaceProvider: localStore,
//		BookingProvider:   localStore,
//		OfferingProvider:  localStore,
//		UserProvider:  localStore,
//		FloorProvider:  localStore,
//	}
//}
