package mail

import (
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

type sendUser struct {
	Name  string
	Email string
}

type ConfirmationParams struct {
	*sendUser
	WorkspaceName string
	FloorName     string
	Start         string
	End           string
}

type EmailClient interface {
	SendConfirmation(typeS string, params *ConfirmationParams) error
	SendCancellation(typeS string, params *ConfirmationParams) error
}

type SendGridClient struct {
	client *sendgrid.Client
}

const Booking = "booking"
const Offering = "offering"
const IWorkUserName = "IWork"
const IWorkEmail = "cs319.icbc@outlook.com"
const EmailBody = `You can view/manage your bookings and offerings under the manage tab at http://icbc-iwork-staging.herokuapp.com/`

func (c *SendGridClient) SendConfirmation(typeS string, params *ConfirmationParams) error {
	from := mail.NewEmail(IWorkUserName, IWorkEmail)
	subject := fmt.Sprintf("Workspace %s confirmed", typeS)
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

func (c *SendGridClient) SendCancellation(typeS string, params *ConfirmationParams) error {
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
