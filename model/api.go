package model

import "time"

type UserAssignment struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Department  string `json:"department"`
	IsAdmin     bool   `json:"is_admin"`
	Email       string `json:"email"`
	WorkspaceId string `json:"workspace_id"`
}

type CreateWorkspaceInput struct {
	WorkspaceId   string `json:"id"`
	WorkspaceName string `json:"name"`
	Props         Attrs  `json:"properties"`
	Details       string `json:"details"`
	UserId        string `json:"user_id"`
}

type BulkCreateWorkspacesInput struct {
	FloorId    string                  `json:"floor_id"`
	Workspaces []*CreateWorkspaceInput `json:"workspaces"`
}

type BulkCreateWorkspaceError struct {
	WorkspaceName string `json:"workspace_name"`
	Message       string `json:"message"`
}

type DeleteFloor struct {
	ForceDelete bool `json:"force_delete"`
}

type LockWorkspaceInput struct {
	WorkspaceId string    `json:"workspace_id"`
	StartDate   time.Time `json:"start_time"`
	EndDate     time.Time `json:"end_time"`
}
