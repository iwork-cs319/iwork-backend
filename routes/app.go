package routes

import (
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go-api/db"
	"go-api/db/postgres"
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
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		log.Println("Failed to connect to redis")
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
}

func (app *App) Close() {
	app.store.Close()
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"1","Hello World!",}`))
}
