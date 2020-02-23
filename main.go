package main

import (
	"go-api/routes"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("Did not find env var PORT, defaulting to 8080")
	}
	dbUrl := os.Getenv("DATABASE_URL")
	gDriveCredentials := os.Getenv("G_DRIVE_CREDENTIALS")
	app := routes.NewApp(dbUrl, gDriveCredentials)
	defer app.Close()
	err := app.Setup(port)
	log.Fatal(err)
}
