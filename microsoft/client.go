package microsoft

import (
	"encoding/json"
	"errors"
	"fmt"
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

const TokenUrl = "https://login.microsoftonline.com/de28de2e-eaf8-4937-a102-735e764a6e31/oauth2/v2.0/token"
const GraphUrl = "https://graph.microsoft.com/v1.0"
const EmailBody = `
	{
	  "message": {
		"subject": "%s",
		"body": {
		  "contentType": "Text",
		  "content": "%s"
		},
		"toRecipients": [
		  {
			"emailAddress": {
			  "address": "%s"
			}
		  }
		]
	  },
	  "saveToSentItems": "true"
	}
`

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
	req, err := http.NewRequest("GET", url, nil)
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
	return c.defaultClient.Do(req)
}

func (c *ADClient) GetAllUsers() (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("%s/users", GraphUrl)
	resp, err := c.doRequest("GET", reqUrl, nil)
	defer resp.Body.Close()
	reqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var body map[string]interface{}
	err = json.Unmarshal(reqBody, &body)
	return body, err
}

func (c *ADClient) SendMail(subject, toEmail, emailBody string) (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("%s/users/%s/sendmail", GraphUrl, c.adminUserId)
	body := strings.NewReader(fmt.Sprintf(EmailBody, subject, emailBody, toEmail))
	req, err := http.NewRequest("POST", reqUrl, body)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	resp, err := client.Do(req)
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
