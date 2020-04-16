package microsoft

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-api/mail"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ADClient struct {
	clientId      string
	scope         string
	clientSecret  string
	accessToken   string
	adminUserId   string
	defaultClient *http.Client
}

func (c *ADClient) SendConfirmation(typeS string, params *mail.EmailParams) error {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return err
	}
	inviteContent := fmt.Sprintf(
		`Your %s for workspace <strong>%s</strong> on floor <strong>%s</strong>
					for the duration of <strong>%s</strong> to <strong>%s</strong> has now been confirmed. 
					<br>%s
					`,
		typeS, params.WorkspaceName,
		params.FloorName,
		params.Start.In(loc).Format("Monday 02 Jan 06 15:04"),
		params.End.In(loc).Format("Monday 02 Jan 06 15:04"),
		mail.EmailBody,
	)
	if typeS == mail.Booking {
		inviteContent = fmt.Sprintf(`%s<br> <a href="%s">Google Map Link</a>`,
			inviteContent,
			fmt.Sprintf(
				"https://www.google.com/maps/search/?api=1&query=%s",
				url.PathEscape(params.FloorAddress),
			),
		)
	}

	return c.sendCalendarInvite(&CalendarInvite{
		subject:   fmt.Sprintf("%s for %s at %s", typeS, params.WorkspaceName, params.FloorName),
		content:   inviteContent,
		startTime: params.Start,
		endTime:   params.End,
		location:  params.WorkspaceName,
		attendees: []*Attendee{
			{
				email: params.Email,
				name:  params.Name,
			},
		},
	})
}

func (c *ADClient) SendCancellation(typeS string, params *mail.EmailParams) error {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return err
	}
	cancellationContent := fmt.Sprintf(
		"Your %s for workspace <strong>%s</strong> on floor <strong>%s</strong> for the duration of <strong>%s</strong> to <strong>%s</strong> has now been cancelled. \n%s",
		typeS, params.WorkspaceName,
		params.FloorName,
		params.Start.In(loc).Format("Monday 02 Jan 06 15:04"),
		params.End.In(loc).Format("Monday 02 Jan 06 15:04"),
		mail.EmailBody,
	)
	return c.sendEmail(&EmailBody{
		subject: fmt.Sprintf("%s cancellation for %s", typeS, params.WorkspaceName),
		content: cancellationContent,
		attendees: []*Attendee{
			{
				email: params.Email,
				name:  params.Name,
			},
		},
	})
}

type Attendee struct {
	email string
	name  string
}

type CalendarInvite struct {
	subject   string
	content   string
	startTime time.Time
	endTime   time.Time
	location  string
	attendees []*Attendee
}

type EmailBody struct {
	subject   string
	content   string
	attendees []*Attendee
}

const TokenUrl = "https://login.microsoftonline.com/de28de2e-eaf8-4937-a102-735e764a6e31/oauth2/v2.0/token"
const GraphUrl = "https://graph.microsoft.com/v1.0"

func buildCalendarInviteBody(invite *CalendarInvite) []byte {
	//{
	//	"attendees": [
	//		{
	//			"emailAddress": {
	//				"address": "bruce@cs319iwork.onmicrosoft.com",
	//				"name": "Bruce Wayne"
	//			},
	//			"type": "required"
	//		}
	//	],
	//	"body": {
	//		"content": "sad",
	//		"contentType": "HTML"
	//	},
	//	"end": {
	//		"dateTime": "2020-03-23T00:03:16",
	//		"timeZone": "Pacific Standard Time"
	//	},
	//	"location": {
	//		"displayName": "W-001"
	//	},
	//	"start": {
	//		"dateTime": "2020-03-21T20:16:36",
	//		"timeZone": "Pacific Standard Time"
	//	},
	//	"subject": "Booking confirmation for W-001 at West 2nd Avenue"
	//}
	body := make(map[string]interface{})
	body["subject"] = invite.subject
	body["body"] = map[string]interface{}{
		"contentType": "HTML",
		"content":     invite.content,
	}
	body["start"] = map[string]interface{}{
		"dateTime": invite.startTime.Format("2006-01-02T15:04:05"),
		"timeZone": "Pacific Standard Time",
	}
	body["end"] = map[string]interface{}{
		"dateTime": invite.endTime.Format("2006-01-02T15:04:05"),
		"timeZone": "Pacific Standard Time",
	}
	body["location"] = map[string]interface{}{
		"displayName": invite.location,
	}
	var attendees []interface{}
	for _, a := range invite.attendees {
		var attendee = make(map[string]interface{})
		attendee["emailAddress"] = map[string]interface{}{
			"address": a.email,
			"name":    a.name,
		}
		attendee["type"] = "required"
		attendees = append(attendees, attendee)
	}
	body["attendees"] = attendees
	res, err := json.Marshal(body)
	if err != nil {
		return nil
	}
	return res
}

func buildEmailBody(email *EmailBody) []byte {
	//{
	//	"message": {
	//		"subject": "Meet for lunch Again?",
	//		"body": {
	//			"contentType": "Text",
	//			"content": "The new cafeteria is open."
	//		},
	//		"toRecipients": [
	//			{
	//				"emailAddress": {
	//					"address": "bruce@cs319iwork.onmicrosoft.com"
	//				}
	//			}
	//		]
	//	},
	//	"saveToSentItems": "true"
	//}
	var attendees []interface{}
	for _, a := range email.attendees {
		var attendee = make(map[string]interface{})
		attendee["emailAddress"] = map[string]interface{}{
			"address": a.email,
			"name":    a.name,
		}
		attendees = append(attendees, attendee)
	}
	body := make(map[string]interface{})
	body["message"] = map[string]interface{}{
		"subject": email.subject,
		"body": map[string]interface{}{
			"contentType": "Text",
			"content":     email.content,
		},
		"toRecipients": attendees,
	}
	body["saveToSentItems"] = true

	res, err := json.Marshal(body)
	if err != nil {
		return nil
	}
	return res
}

func NewADClient(clientId, scope, clientSecret, adminUserId string) (*ADClient, error) {
	client := &ADClient{
		clientId:     clientId,
		scope:        scope,
		clientSecret: clientSecret,
		adminUserId:  adminUserId,
		defaultClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	err := client.RefreshToken()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *ADClient) RefreshToken() error {
	if err := c.Ping(); err == nil {
		return nil
	}
	data := url.Values{}
	data.Set("client_id", c.clientId)
	data.Set("scope", c.scope)
	data.Set("client_secret", c.clientSecret)
	data.Set("grant_type", "client_credentials")

	resp, err := http.Post(TokenUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("incorrect response code while fetching token, received %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	reqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var body map[string]interface{}
	err = json.Unmarshal(reqBody, &body)
	accessToken := fmt.Sprintf("%v", body["access_token"])
	if accessToken == "" {
		return errors.New("no access token found")
	}
	c.accessToken = accessToken
	return nil
}

func (c *ADClient) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	return req, nil
}

func (c *ADClient) doRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.defaultClient.Do(req)
}

func (c *ADClient) Ping() error {
	resp, err := c.doRequest("GET", GraphUrl, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("token likely expired; statusCode=%d", resp.StatusCode))
	}
	return nil
}

func (c *ADClient) GetAllUsers() (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("%s/users", GraphUrl)
	resp, err := c.doRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	reqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var body map[string]interface{}
	err = json.Unmarshal(reqBody, &body)
	return body, err
}

func (c *ADClient) sendCalendarInvite(invite *CalendarInvite) error {
	_ = c.RefreshToken()
	reqUrl := fmt.Sprintf("%s/users/%s/events", GraphUrl, c.adminUserId)
	body := bytes.NewBuffer(buildCalendarInviteBody(invite))
	resp, err := c.doRequest("POST", reqUrl, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("failed to send email, %d", resp.StatusCode))
	}
	return nil
}

func (c *ADClient) sendEmail(email *EmailBody) error {
	_ = c.RefreshToken()
	reqUrl := fmt.Sprintf("%s/users/%s/sendmail", GraphUrl, c.adminUserId)
	body := bytes.NewBuffer(buildEmailBody(email))
	resp, err := c.doRequest("POST", reqUrl, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("failed to send email, %d", resp.StatusCode))
	}
	return nil
}
