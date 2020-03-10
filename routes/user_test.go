package routes

import (
	"encoding/json"
	"go-api/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testUsersEndpoints(t *testing.T, app *App) {

}

func testGetAllUsers(t *testing.T, app *App) {
	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.GetAllUsers)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var payload []*model.User
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	if len(payload) != 4 {
		t.Fa
	}
}
