package routes

import (
	"github.com/gorilla/mux"
	"go-api/db"
	"net/http"
)

type App struct {
	router *mux.Router
	store  *db.DataStore
}

func NewApp() *App {
	return &App{
		router: mux.NewRouter().StrictSlash(true),
		store:  db.NewLocalDataStore(),
	}
}

func (app *App) Setup(port string) error {
	app.router.HandleFunc("/", app.index)
	app.RegisterWorkspaceRoutes()
	return http.ListenAndServe(":"+port, app.router)
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"1","Hello World!",}`))
}
