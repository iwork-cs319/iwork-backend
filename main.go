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
	err := routes.NewApp().Setup(port)
	log.Fatal(err)
}