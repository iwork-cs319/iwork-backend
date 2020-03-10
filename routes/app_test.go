package routes

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"go-api/db/postgres"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type testRouteConfig struct {
	Method    string
	URL       string
	Body      io.Reader
	URLParams map[string]string
	Handler   func(http.ResponseWriter, *http.Request)
}

func executeReq(t *testing.T, config *testRouteConfig) *httptest.ResponseRecorder {
	req, err := http.NewRequest(config.Method, config.URL, config.Body)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, config.URLParams)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(config.Handler)
	handler.ServeHTTP(rr, req)
	return rr
}

type AppTestSuite struct {
	suite.Suite
	app *App
}

func (suite *AppTestSuite) SetupTest() {
	suite.app = NewTestApp()
}

func TestApp(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func NewTestApp() *App {
	dbUrl := os.Getenv("TEST_DB_URL")
	store, err := postgres.NewPostgresDataStore(dbUrl)
	if err != nil {
		log.Println("Failed to connect to database")
		log.Fatal(err)
	}
	//gDriveConfig := os.Getenv("G_DRIVE_CREDENTIALS")
	//driveClient, err := db.NewDriveClient(gDriveConfig)
	//if err != nil {
	//	log.Println("Failed to connect to google drive")
	//	log.Fatal(err)
	//}
	return &App{
		router: mux.NewRouter().StrictSlash(true),
		store:  store,
		gDrive: nil,
	}
}
