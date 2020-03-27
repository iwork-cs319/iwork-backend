package model

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
