package main

import (
	"go-api/mail"
	"go-api/microsoft"
	"log"
	"time"
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
		"380ec12d-9d53-45a6-91da-c3a262cb3dca",
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("%v", users)
	err = c.SendConfirmation("booking", &mail.EmailParams{
		Name:          "Bruce Wayne",
		Email:         "bruce@cs319iwork.onmicrosoft.com",
		WorkspaceName: "W-001",
		FloorName:     "West 2nd Avenue",
		Start:         time.Unix(1584846996, 0),
		End:           time.Unix(1584946996, 0),
	})
	log.Println(err)
}
