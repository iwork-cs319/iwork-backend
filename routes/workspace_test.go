package routes

import (
	"encoding/json"
	"fmt"
	"go-api/model"
	"net/http"
	"testing"
)

var Workspace1 = &model.Workspace{
	ID:    "a12411d3-d281-3735-b000-bf94b094d2af",
	Name:  "W-001",
	Floor: MainFloor.ID,
	Props: nil,
}
var Workspace2 = &model.Workspace{
	ID:    "f2188d5e-509f-3086-a031-d86da93a55c4",
	Name:  "W-002",
	Floor: MainFloor.ID,
	Props: nil,
}
var Workspace3 = &model.Workspace{
	ID:    "62eb266e-740d-3973-9e50-77b0287e3026",
	Name:  "W-003",
	Floor: MainFloor.ID,
	Props: nil,
}
var Workspace4 = &model.Workspace{
	ID:    "5e56de3d-2323-372d-897f-23d6037c8581",
	Name:  "W-004",
	Floor: MainFloor.ID,
	Props: nil,
}
var Workspace5 = &model.Workspace{
	ID:    "aad40cbb-4baf-3931-a5d2-6f98b414182a",
	Name:  "W-005",
	Floor: MainFloor.ID,
	Props: nil,
}
var Workspace6 = &model.Workspace{
	ID:    "bb15369d-e6e0-33b8-8b97-1779f8865890",
	Name:  "W-006",
	Floor: MainFloor.ID,
	Props: nil,
}
var Workspace7 = &model.Workspace{
	ID:    "3361d373-781a-34d7-bbb8-c7d562a0cf51",
	Name:  "W-007",
	Floor: MainFloor.ID,
	Props: nil,
}

func testWorkspaceEndpoints(t *testing.T, app *App) {
	testGetAvailable(t, app)
	testGetOneWorkspace(t, app)
	testGetOneWorkspaceFail(t, app)
	testGetAllWorkspaces(t, app)
}

func testGetAvailable(t *testing.T, app *App) {
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: app.GetAvailability,
		URL: fmt.Sprintf(
			"/workspaces/available?floor=%s&start=%s&end=%s",
			MainFloor.ID,
			"1547337600",
			"1547596800",
		),
	})
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("testGetAvailable: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var payload []*string
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	if len(payload) != 2 {
		t.Fatalf("testGetAvailable: expected 2 workspaces got %d", len(payload))
	}

}

func testGetOneWorkspace(t *testing.T, app *App) {
	workspaceId := Workspace1.ID
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: app.GetOneWorkspace,
		URL:     fmt.Sprintf("/workspaces/%s", workspaceId),
		URLParams: map[string]string{
			"id": workspaceId,
		},
	})
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("testGetOneWorkspace: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var payload *model.Workspace
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	if Workspace1.Equal(payload) {
		t.Fatalf("testGetOneWorkspace: incorrect workspace : got %s", payload.Name)
	}
}

func testGetOneWorkspaceFail(t *testing.T, app *App) {
	workspaceId := "2"
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: app.GetOneWorkspace,
		URL:     fmt.Sprintf("/workspaces/%s", workspaceId),
		URLParams: map[string]string{
			"id": workspaceId,
		},
	})
	if status := rr.Code; status != http.StatusNotFound {
		t.Fatalf("testGetOneWorkspaceFail: handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func testGetAllWorkspaces(t *testing.T, app *App) {
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: app.GetAllWorkspaces,
		URL:     "/workspaces/",
	})

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("testGetAllWorkspace: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var payload []*model.Workspace
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	if len(payload) != 7 {
		t.Fatalf("testGetAllWorkspace: incorrect number of workspaces: got %d users", len(payload))
	}
	expectedWorkspaces := []*model.Workspace{Workspace1, Workspace2, Workspace3, Workspace4, Workspace5, Workspace6, Workspace7}
	for _, expected := range expectedWorkspaces {
		found := false
		for _, workspace := range payload {
			if expected.Equal(workspace) {
				found = true
			}
		}
		if !found {
			t.Fatalf("testGetAllWorkspace: %s not found in user list", expected)
		}
	}
}
