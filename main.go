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
	msClientId := os.Getenv("MICROSOFT_CLIENT_ID")
	msGraphScope := os.Getenv("MICROSOFT_GRAPH_SCOPE")
	msClientSecret := os.Getenv("MICROSOFT_CLIENT_SECRET")
	adminUserId := os.Getenv("ADMIN_USER_ID")
	app := routes.NewApp(&routes.AppConfig{
		DbUrl:          dbUrl,
		GDriveConfig:   gDriveCredentials,
		MsClientId:     msClientId,
		MsScope:        msGraphScope,
		MsClientSecret: msClientSecret,
		AdminUserId:    adminUserId,
	})
	defer app.Close()
	err := app.Setup(port)
	log.Fatal(err)
}
