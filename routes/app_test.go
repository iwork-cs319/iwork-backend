package routes

import (
	"github.com/gorilla/mux"
	"go-api/db"
	"go-api/db/postgres"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type testRouteConfig struct {
	Method  string
	URL     string
	Path    string
	Body    io.Reader
	Handler func(http.ResponseWriter, *http.Request)
}

func doRequest(t *testing.T, config *testRouteConfig) *httptest.ResponseRecorder {
	req, err := http.NewRequest(config.Method, config.URL, config.Body)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(config.Path, config.Handler)
	router.ServeHTTP(rr, req)
	return rr
}

func NewTestApp() *App {
	dbUrl := os.Getenv("TEST_DB_URL")
	store, err := postgres.NewPostgresDataStore(dbUrl)
	if err != nil {
		log.Println("Failed to connect to database")
		log.Fatal(err)
	}
	gDriveConfig := os.Getenv("G_DRIVE_CREDENTIALS")
	driveClient, err := db.NewDriveClient(gDriveConfig)
	if err != nil {
		log.Println("Failed to connect to google drive")
		log.Fatal(err)
	}
	return &App{
		router: mux.NewRouter().StrictSlash(true),
		store:  store,
		gDrive: driveClient,
	}
}

func TestApp(t *testing.T) {
	a := NewTestApp()
	testUsersEndpoints(t, a)
	//TODO add more here
}
