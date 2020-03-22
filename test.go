package main

import (
	"go-api/microsoft"
	"log"
)

func main() {
	//c, err := microsoft.NewADClient(
	//	os.Getenv("MICROSOFT_CLIENT_ID"),
	//	os.Getenv("MICROSOFT_GRAPH_SCOPE"),
	//	os.Getenv("MICROSOFT_CLIENT_SECRET"),
	//)
	c, err := microsoft.NewADClient(
		"713700f2-6fcf-4258-94bc-0c7986180460",
		"https://graph.microsoft.com/.default",
		"IL2qOE:m_P/pJ?AiGPDqO10MeoJzZp23",
	)
	if err != nil {
		log.Fatal(err)
	}
	users, err := c.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", users)
}
