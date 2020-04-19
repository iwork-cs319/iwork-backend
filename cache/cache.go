package cache

import "time"

type Cache interface {
	CheckWorkspaceLock(workspaceId string, start, end time.Time) (bool, error)
	CreateWorkspaceLock(workspaceId string, start, end time.Time) (string, error)
	DeleteWorkspaceLock(workspaceId string, key string) error
}
