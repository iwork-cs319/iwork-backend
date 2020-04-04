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
	"time"
)

type App struct {
	router *mux.Router
	store  *db.DataStore
	gDrive db.Drive
	cache  *redis.Pool
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
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
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
	redisCache := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisUrl)
		},
	}
	return &App{
		router: mux.NewRouter().StrictSlash(true),
		store:  store,
		gDrive: driveClient,
		email:  msClient,
		cache:  redisCache,
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
	app.RegisterUserRoutes()
	app.RegisterFloorRoutes()
	app.RegisterWorkspaceRoutes()
	app.RegisterBookingRoutes()
	app.RegisterOfferingRoutes()
	app.RegisterArchiverRoutes()
}

func (app *App) Close() {
	app.store.Close()
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"1","Hello World!",}`))
}
