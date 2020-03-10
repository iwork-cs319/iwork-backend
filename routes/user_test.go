package routes

import (
	"encoding/json"
	"fmt"
	"go-api/model"
	"net/http"
	"testing"
)

// main entry point for user endpoints' tests
func testUsersEndpoints(t *testing.T, app *App) {
	testGetAllUsers(t, app)
	testGetOneUser(t, app)
	testGetOneUserFail(t, app)
}

func testGetOneUserFail(t *testing.T, app *App) {
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: app.GetOneUser,
		URL:     fmt.Sprintf("/users/%s", "2"),
		URLParams: map[string]string{
			"id": "2",
		},
	})

	if status := rr.Code; status != http.StatusNotFound {
		t.Fatalf("testGetOneUserFail: handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func testGetOneUser(t *testing.T, app *App) {
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: app.GetOneUser,
		URL:     fmt.Sprintf("/users/%s", "e99a988a-1d41-3997-8d59-959a48ac24a0"),
		URLParams: map[string]string{
			"id": "e99a988a-1d41-3997-8d59-959a48ac24a0",
		},
	})

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("testGetOneUser: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var payload *model.User
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	if !payload.IsAdmin && payload.Name != "Barry Allen" {
		t.Fatalf("testGetOneUser: incorrect user : got %s", payload.Name)
	}
}

func testGetAllUsers(t *testing.T, app *App) {
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: app.GetAllUsers,
		URL:     "/users/",
	})

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("testGetAllUsers: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var payload []*model.User
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	if len(payload) != 4 {
		t.Fatalf("testGetAllUsers: incorrect number of users: got %d users", len(payload))
	}
	expectedUsers := []string{
		"Barry Allen",
		"Bruce Wayne",
		"Clark Kent",
		"Diana Prince",
	}
	for _, uName := range expectedUsers {
		found := false
		for _, u := range payload {
			if u.Name == uName {
				found = true
			}
		}
		if !found {
			t.Fatalf("testGetAllUsers: %s not found in user list", uName)
		}
	}
}
