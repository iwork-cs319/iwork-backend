package model

import "time"

type Workspace struct {
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Floor string                 `json:"floor_id"`
	Props map[string]interface{} `json:"props"`
}

func (w *Workspace) Equal(other *Workspace) bool {
	return w.ID == other.ID && w.Name == other.Name && w.Floor == other.Floor
}

type Booking struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	StartDate   time.Time `json:"start_time"`
	EndDate     time.Time `json:"end_time"`
	Cancelled   bool      `json:"cancelled"`
	CreatedBy   string    `json:"created_by"`
}

type ExpandedBooking struct {
	Booking
	WorkspaceName string `json:"workspace_name"`
	UserName      string `json:"user_name"`
	FloorID       string `json:"floor_id"`
	FloorName     string `json:"floor_name"`
}

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
	IsAdmin    bool   `json:"is_admin"`
	Email      string `json:"email"`
}

func (w *User) Equal(other *User) bool {
	return w.ID == other.ID && w.Name == other.Name &&
		w.Email == other.Email && w.Department == other.Department &&
		w.IsAdmin == other.IsAdmin
}

type Floor struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

type Offering struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	StartDate   time.Time `json:"start_time"`
	EndDate     time.Time `json:"end_time"`
	Cancelled   bool      `json:"cancelled"`
	CreatedBy   string    `json:"created_by"`
}

type ExpandedOffering struct {
	Offering
	WorkspaceName string `json:"workspace_name"`
	UserName      string `json:"user_name"`
	FloorID       string `json:"floor_id"`
	FloorName     string `json:"floor_name"`
}
