package model

type Workspace struct {
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Floor string                 `json:"floor"`
	Props map[string]interface{} `json:"props"`
}
