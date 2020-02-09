package routes

import (
	"github.com/gorilla/mux"
	"go-api/db"
	"log"
	"net/http"
)

type App struct {
	router *mux.Router
	store  *db.DataStore
}

func NewApp(dbUrl string) *App {
	store, err := db.NewPostgresDataStore(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	return &App{
		router: mux.NewRouter().StrictSlash(true),
		store:  store,
	}
}

func (app *App) Setup(port string) error {
	app.router.HandleFunc("/", app.index)
	app.RegisterWorkspaceRoutes()
	app.RegisterBookingRoutes()
	log.Println("App running at port:", port)
	return http.ListenAndServe(":"+port, app.router)
}

func (app *App) Close() {
	app.store.Close()
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"1","Hello World!",}`))
}
