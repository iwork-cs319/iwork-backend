package mail

import (
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

type SendGridClient struct {
	client *sendgrid.Client
}

const Booking = "booking"
const Offering = "offering"
const IWorkUserName = "IWork"
const IWorkEmail = "cs319.icbc@outlook.com"
const EmailBody = `
	You can manage your bookings and offerings under the manage tab at <a href="http://icbc-iwork-staging.herokuapp.com/">http://icbc-iwork-staging.herokuapp.com/</a>
	Please note that cancelling this invite wont cancel this action. Please contact an admin in this scenario
`

func (c *SendGridClient) SendConfirmation(typeS string, params *EmailParams) error {
	from := mail.NewEmail(IWorkUserName, IWorkEmail)
	subject := fmt.Sprintf("Workspace %s confirmed", typeS)
	to := mail.NewEmail(params.Name, IWorkEmail)
	plainTextContent := fmt.Sprintf(
		"Your %s for workspace %s on floor %s for the duration of %s to %s has now been confirmed. \n%s",
		typeS, params.WorkspaceName, params.FloorName, params.Start, params.End, EmailBody,
	)
	htmlContent := plainTextContent
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	_, err := c.client.Send(message)
	if err != nil {
		log.Printf("SendGrid.SendConfirmation: failed to send email: %+v", err)
		return err
	}
	return nil
}

func (c *SendGridClient) SendCancellation(typeS string, params *EmailParams) error {
	from := mail.NewEmail(IWorkUserName, IWorkEmail)
	subject := fmt.Sprintf("Workspace %s cancelled", typeS)
	to := mail.NewEmail(params.Name, params.Email)
	plainTextContent := fmt.Sprintf(
		"Your %s for workspace %s on floor %s for the duration of %s to %s has now been confirmed. \n%s",
		typeS, params.WorkspaceName, params.FloorName, params.Start, params.End, EmailBody,
	)
	htmlContent := plainTextContent
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	_, err := c.client.Send(message)
	if err != nil {
		log.Printf("SendGrid.SendConfirmation: failed to send email: %+v", err)
		return err
	}
	return nil
}

func NewSendGridClient() (EmailClient, error) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return nil, errors.New("API Key missing")
	}
	return &SendGridClient{
		client: sendgrid.NewSendClient(apiKey),
	}, nil
}
