package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-api/model"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

var UserDefault = &model.User{
	ID:         "decade00-0000-4000-a000-000000000000",
	Name:       "Default User",
	Department: "N/A",
	IsAdmin:    false,
	Email:      "N/A",
}
var UserBarry = &model.User{
	ID:         "e99a988a-1d41-3997-8d59-959a48ac24a0",
	Name:       "Barry Allen",
	Department: "R&D",
	IsAdmin:    false,
	Email:      "",
}
var UserBruce = &model.User{
	ID:         "8b5bb736-6a1d-3378-8e71-ab45fe8beb84",
	Name:       "Bruce Wayne",
	Department: "Engineering",
	IsAdmin:    true,
	Email:      "",
}
var UserClark = &model.User{
	ID:         "32ea2fb1-7124-304a-b9c3-eb445578103e",
	Name:       "Clark Kent",
	Department: "Marketing",
	IsAdmin:    false,
	Email:      "",
}
var UserDiana = &model.User{
	ID:         "ab6c2c4f-c112-3c1e-bcf1-42cdea289c1f",
	Name:       "Diana Prince",
	Department: "Operations",
	IsAdmin:    false,
	Email:      "",
}

func (suite *AppTestSuite) TestGetOneUserFail() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneUser,
		URL:     fmt.Sprintf("/users/%s", "2"),
		URLParams: map[string]string{
			"id": "2",
		},
	})

	assert.Equal(t, http.StatusNotFound, rr.Code, "status code")
}

func (suite *AppTestSuite) TestGetOneUser() {
	userId := UserBarry.ID
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneUser,
		URL:     fmt.Sprintf("/users/%s", userId),
		URLParams: map[string]string{
			"id": userId,
		},
	})

	require.Equal(t, rr.Code, http.StatusOK, "status code")
	var payload *model.User
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, UserBarry, payload, "incorrect response object")
}

func (suite *AppTestSuite) TestGetAllUsers() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetAllUsers,
		URL:     "/users/",
	})

	require.Equal(t, rr.Code, http.StatusOK, "status code")

	var payload []*model.User
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 5, len(payload), "incorrect response size")
	assert.Contains(t, payload, UserDefault, "doesnt contain defaultUsers")
	assert.Contains(t, payload, UserBarry, "doesnt contain barry")
	assert.Contains(t, payload, UserDiana, "doesnt contain diana")
	assert.Contains(t, payload, UserClark, "doesnt contain clark")
	assert.Contains(t, payload, UserBruce, "doesnt contain bruce")
}

// Run after
func (suite *AppTestSuite) Test_CreateUsers() {
	t := suite.T()
	body := new(bytes.Buffer)
	file, err := os.Open("../test-fixtures/users.csv")
	if err != nil {
		t.Fatalf("failed to open file")
	}
	fileContents, _ := ioutil.ReadAll(file)

	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("users", "users.csv")
	_, _ = part.Write(fileContents)
	_ = writer.Close()

	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    body,
		Handler: suite.app.CreateUsers,
		URL:     "/users/",
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	})
	if !assert.Equal(t, http.StatusCreated, rr.Code, "status code") {
		return
	}

	var payload []*model.User
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 1, len(payload), "incorrect response size")
	user := payload[0]
	assert.Equal(t, "3f9fc7c0-d675-40be-9ad1-9babfad625d7", user.ID)
	assert.Equal(t, "fake_user@iworkcs319.onmicrosoft.com", user.Email)
	assert.Equal(t, "Fake User", user.Name)
	assert.Equal(t, "Operations", user.Department)
	assert.False(t, user.IsAdmin)
}

func (suite *AppTestSuite) Test_CreateUsersEmptyFile() {
	t := suite.T()
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    body,
		Handler: suite.app.CreateUsers,
		URL:     "/users/",
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	})
	if !assert.Equal(t, http.StatusBadRequest, rr.Code, "status code") {
		return
	}
}

func (suite *AppTestSuite) Test_CreateUsersBadContentType() {
	t := suite.T()
	body := new(bytes.Buffer)

	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    body,
		Handler: suite.app.CreateUsers,
		URL:     "/users/",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if !assert.Equal(t, http.StatusBadRequest, rr.Code, "status code") {
		return
	}
}

func (suite *AppTestSuite) Test_CreateUsersBadFile() {
	t := suite.T()
	body := new(bytes.Buffer)
	file, err := os.Open("../test-fixtures/workspaces.csv")
	if err != nil {
		t.Fatalf("failed to open file")
	}
	fileContents, _ := ioutil.ReadAll(file)

	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("users", "users.csv")
	_, _ = part.Write(fileContents)
	_ = writer.Close()

	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    body,
		Handler: suite.app.CreateUsers,
		URL:     "/users/",
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	})
	if !assert.Equal(t, http.StatusBadRequest, rr.Code, "status code") {
		return
	}
}

func (suite *AppTestSuite) Test_CreateUsersBadFile2() {
	t := suite.T()
	body := new(bytes.Buffer)
	file, err := os.Open("../test-fixtures/test-img.jpg")
	if err != nil {
		t.Fatalf("failed to open file")
	}
	fileContents, _ := ioutil.ReadAll(file)

	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("users", "users.csv")
	_, _ = part.Write(fileContents)
	_ = writer.Close()

	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    body,
		Handler: suite.app.CreateUsers,
		URL:     "/users/",
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	})
	if !assert.Equal(t, http.StatusBadRequest, rr.Code, "status code") {
		return
	}
}
