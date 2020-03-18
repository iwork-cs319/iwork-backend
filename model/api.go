package model

type UserAssignment struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Department  string `json:"department"`
	IsAdmin     bool   `json:"is_admin"`
	Email       string `json:"email"`
	WorkspaceId string `json:"workspace_id"`
}
