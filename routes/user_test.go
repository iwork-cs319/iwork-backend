package routes

import (
	"encoding/json"
	"fmt"
	"go-api/model"
	"net/http"
	"testing"
)

var UserBarry = &model.User{
	ID:         "e99a988a-1d41-3997-8d59-959a48ac24a0",
	Name:       "Barry Allen",
	Department: "R&D",
	IsAdmin:    false,
	Email:      "",
}
var UserBruce = &model.User{
	ID:         "8b5bb736-6a1d-3378-8e71-ab45fe8beb84",
	Name:       "Bruce Wayne",
	Department: "Engineering",
	IsAdmin:    true,
	Email:      "",
}
var UserClark = &model.User{
	ID:         "32ea2fb1-7124-304a-b9c3-eb445578103e",
	Name:       "Clark Kent",
	Department: "Marketing",
	IsAdmin:    false,
	Email:      "",
}
var UserDiana = &model.User{
	ID:         "ab6c2c4f-c112-3c1e-bcf1-42cdea289c1f",
	Name:       "Diana Prince",
	Department: "Operations",
	IsAdmin:    false,
	Email:      "",
}

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
	userId := UserBarry.ID
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: app.GetOneUser,
		URL:     fmt.Sprintf("/users/%s", userId),
		URLParams: map[string]string{
			"id": userId,
		},
	})

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("testGetOneUser: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var payload *model.User
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	if UserBarry.Equal(payload) {
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

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("testGetAllUsers: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var payload []*model.User
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	if len(payload) != 4 {
		t.Fatalf("testGetAllUsers: incorrect number of users: got %d users", len(payload))
	}
	expectedUsers := []*model.User{UserBarry, UserBruce, UserClark, UserDiana}
	for _, expected := range expectedUsers {
		found := false
		for _, u := range payload {
			if u.Equal(expected) {
				found = true
			}
		}
		if !found {
			t.Fatalf("testGetAllUsers: %s not found in user list", expected.Name)
		}
	}
}
