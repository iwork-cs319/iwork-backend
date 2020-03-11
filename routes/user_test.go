package routes

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-api/model"
	"net/http"
)

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

	assert.Equal(t, rr.Code, http.StatusOK, "status code")
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

	assert.Equal(t, rr.Code, http.StatusOK, "status code")

	var payload []*model.User
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 4, len(payload), "incorrect response size")
	assert.Contains(t, payload, UserBarry, "doesnt contain barry")
	assert.Contains(t, payload, UserDiana, "doesnt contain diana")
	assert.Contains(t, payload, UserClark, "doesnt contain clark")
	assert.Contains(t, payload, UserBruce, "doesnt contain bruce")
}
