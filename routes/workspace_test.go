package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func testWorkspaceEndpoints(t *testing.T, app *App) {
	testGetAvailable(t, app)
}

func testGetAvailable(t *testing.T, app *App) {
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: app.GetAvailability,
		URL: fmt.Sprintf(
			"/workspaces/available?floor=%s&start=%s&end=%s",
			"a709a01e-e9e5-3139-8f0d-fba4c0a2187f",
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
