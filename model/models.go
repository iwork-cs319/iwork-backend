package model

import "time"

type Workspace struct {
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Floor string                 `json:"floor_id"`
	Props map[string]interface{} `json:"props"`
}

type Booking struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	StartDate   time.Time `json:"start_time"`
	EndDate     time.Time `json:"end_time"`
	Cancelled   bool      `json:"cancelled"`
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
}

type ExpandedOffering struct {
	Offering
	WorkspaceName string `json:"workspace_name"`
	UserName      string `json:"user_name"`
	FloorID       string `json:"floor_id"`
	FloorName     string `json:"floor_name"`
}
