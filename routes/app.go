package routes

import (
	"github.com/gorilla/mux"
	"go-api/db"
	"go-api/db/postgres"
	"log"
	"net/http"
)

type App struct {
	router *mux.Router
	store  *db.DataStore
	gDrive *db.Drive
}

func NewApp(dbUrl, gDriveConfig string) *App {
	store, err := postgres.NewPostgresDataStore(dbUrl)
	if err != nil {
		log.Println("Failed to connect to database")
		log.Fatal(err)
	}
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

func (app *App) Setup(port string) error {
	app.router.StrictSlash(true)
	app.router.HandleFunc("/", app.index)
	app.RegisterRoutes()
	log.Println("App running at port:", port)
	return http.ListenAndServe(":"+port, app.router)
}

func (app *App) RegisterRoutes() {
	app.RegisterWorkspaceRoutes()
	app.RegisterBookingRoutes()
	app.RegisterOfferingRoutes()
	app.RegisterUserRoutes()
	app.RegisterFloorRoutes()
}

func (app *App) Close() {
	app.store.Close()
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"1","Hello World!",}`))
}
