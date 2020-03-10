package routes

import (
	"github.com/gorilla/mux"
	"go-api/db"
	"go-api/db/postgres"
	"log"
	"os"
	"testing"
)

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
}
