package routes

import (
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go-api/db"
	"go-api/db/postgres"
	"go-api/mail"
	"go-api/microsoft"
	"log"
	"net/http"
	"os"
)

type App struct {
	router *mux.Router
	store  *db.DataStore
	gDrive db.Drive
	cache  redis.Conn
	email  mail.EmailClient
}

type AppConfig struct {
	DbUrl          string
	GDriveConfig   string
	MsClientId     string
	MsScope        string
	MsClientSecret string
	AdminUserId    string
}

func NewApp(config *AppConfig) *App {
	store, err := postgres.NewPostgresDataStore(config.DbUrl)
	if err != nil {
		log.Println("Failed to connect to database")
		log.Fatal(err)
	}
	driveClient, err := db.NewDriveClient(config.GDriveConfig)
	if err != nil {
		log.Println("Failed to connect to google drive")
		log.Fatal(err)
	}
	cache, err := redis.DialURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Println("Failed to connect to redis")
		log.Fatal(err)
	}
	msClient, err := microsoft.NewADClient(
		config.MsClientId,
		config.MsScope,
		config.MsClientSecret,
		config.AdminUserId,
	)
	if err != nil {
		log.Println("Failed to create AD Client")
		log.Fatal(err)
	}
	return &App{
		router: mux.NewRouter().StrictSlash(true),
		store:  store,
		gDrive: driveClient,
		cache:  cache,
		email:  msClient,
	}
}

func (app *App) Setup(port string) error {
	app.router.HandleFunc("/", app.index)
	app.RegisterRoutes()
	log.Println("App running at port:", port)
	handler := cors.AllowAll().Handler(app.router)
	return http.ListenAndServe(":"+port, handler)
}

func (app *App) RegisterRoutes() {
	app.RegisterLoginRoutes()
	app.RegisterUserRoutes()
	app.RegisterFloorRoutes()
	app.RegisterWorkspaceRoutes()
	app.RegisterBookingRoutes()
	app.RegisterOfferingRoutes()
}

func (app *App) Close() {
	app.store.Close()
	app.cache.Close()
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"1","Hello World!",}`))
}
