package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-api/model"
	"net/http"
)

var Workspace1 = &model.Workspace{
	ID:      "a12411d3-d281-3735-b000-bf94b094d2af",
	Name:    "W-001",
	Floor:   MainFloor.ID,
	Props:   map[string]interface{}{},
	Details: "",
}
var Workspace2 = &model.Workspace{
	ID:      "f2188d5e-509f-3086-a031-d86da93a55c4",
	Name:    "W-002",
	Floor:   MainFloor.ID,
	Props:   map[string]interface{}{},
	Details: "",
}
var Workspace3 = &model.Workspace{
	ID:      "62eb266e-740d-3973-9e50-77b0287e3026",
	Name:    "W-003",
	Floor:   MainFloor.ID,
	Props:   map[string]interface{}{},
	Details: "",
}
var Workspace4 = &model.Workspace{
	ID:      "5e56de3d-2323-372d-897f-23d6037c8581",
	Name:    "W-004",
	Floor:   MainFloor.ID,
	Props:   map[string]interface{}{},
	Details: "",
}
var Workspace5 = &model.Workspace{
	ID:      "aad40cbb-4baf-3931-a5d2-6f98b414182a",
	Name:    "W-005",
	Floor:   MainFloor.ID,
	Props:   map[string]interface{}{},
	Details: "",
}
var Workspace6 = &model.Workspace{
	ID:      "bb15369d-e6e0-33b8-8b97-1779f8865890",
	Name:    "W-006",
	Floor:   MainFloor.ID,
	Props:   map[string]interface{}{},
	Details: "",
}
var Workspace7 = &model.Workspace{
	ID:      "3361d373-781a-34d7-bbb8-c7d562a0cf51",
	Name:    "W-007",
	Floor:   MainFloor.ID,
	Props:   map[string]interface{}{},
	Details: "",
}

func (suite *AppTestSuite) TestGetAvailable() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetAvailability,
		URL: fmt.Sprintf(
			"/workspaces/available?floor=%s&start=%s&end=%s",
			MainFloor.ID,
			"1547337600",
			"1547596800",
		),
	})
	require.Equal(t, rr.Code, http.StatusOK, "status code")

	var payload []*string
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 2, len(payload), "incorrect response size")
}

func (suite *AppTestSuite) TestGetOneWorkspace() {
	workspaceId := Workspace1.ID
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneWorkspace,
		URL:     fmt.Sprintf("/workspaces/%s", workspaceId),
		URLParams: map[string]string{
			"id": workspaceId,
		},
	})
	require.Equal(t, rr.Code, http.StatusOK, "status code")

	var payload *model.Workspace
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, Workspace1, payload, "incorrect response object")
}

func (suite *AppTestSuite) TestGetOneWorkspaceFail() {
	workspaceId := "2"
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneWorkspace,
		URL:     fmt.Sprintf("/workspaces/%s", workspaceId),
		URLParams: map[string]string{
			"id": workspaceId,
		},
	})
	require.Equal(t, rr.Code, http.StatusNotFound, "status code")
}

func (suite *AppTestSuite) TestGetAllWorkspaces() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetAllWorkspaces,
		URL:     "/workspaces/",
	})

	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusOK, "status code")

	// Check the response body is what we expect.
	var payload []*model.Workspace
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 7, len(payload), "incorrect response size")
	assert.Contains(t, payload, Workspace1, "doesnt contain workspace1")
	assert.Contains(t, payload, Workspace2, "doesnt contain workspace2")
	assert.Contains(t, payload, Workspace3, "doesnt contain workspace3")
	assert.Contains(t, payload, Workspace4, "doesnt contain workspace4")
	assert.Contains(t, payload, Workspace5, "doesnt contain workspace5")
	assert.Contains(t, payload, Workspace6, "doesnt contain workspace6")
	assert.Contains(t, payload, Workspace7, "doesnt contain workspace7")
}

func (suite *AppTestSuite) TestZZCrUDWorkspace() {
	t := suite.T()
	newWorkspace := &model.Workspace{
		Name:  "fake-workspace",
		Floor: MainFloor.ID,
		Props: nil,
	}
	body := new(bytes.Buffer)
	datum, _ := json.Marshal(newWorkspace)
	body.Write(datum)

	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    body,
		Handler: suite.app.CreateWorkspace,
		URL:     "/users/",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if !assert.Equal(t, http.StatusCreated, rr.Code, "failed to create") {
		return
	}
	var payload *model.Workspace
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.NotEmpty(t, payload.ID)
	assert.Equal(t, newWorkspace.Name, payload.Name)
	assert.Equal(t, newWorkspace.Floor, payload.Floor)
	assert.Equal(t, newWorkspace.Props, payload.Props)

	updatedWorkspace := &model.Workspace{
		Name:  "updated-workspace",
		Floor: MainFloor.ID,
		Props: nil,
	}
	body = new(bytes.Buffer)
	datum, _ = json.Marshal(updatedWorkspace)
	body.Write(datum)

	rr = executeReq(t, &testRouteConfig{
		Method:  http.MethodPatch,
		Body:    body,
		Handler: suite.app.UpdateWorkspace,
		URL:     fmt.Sprintf("/users/%s", payload.ID),
		URLParams: map[string]string{
			"id": payload.ID,
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if !assert.Equal(t, http.StatusOK, rr.Code, "failed to update") {
		return
	}
	var payloadUpdate *model.Workspace
	_ = json.Unmarshal(rr.Body.Bytes(), &payloadUpdate)
	assert.Equal(t, payload.ID, payloadUpdate.ID) // check if uuid is the same
	assert.Equal(t, updatedWorkspace.Name, payloadUpdate.Name)
	assert.Equal(t, updatedWorkspace.Floor, payloadUpdate.Floor)
	assert.Equal(t, updatedWorkspace.Props, payloadUpdate.Props)
}
