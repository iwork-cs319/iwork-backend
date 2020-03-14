package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-api/model"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

var MainFloor = &model.Floor{
	ID:          "a709a01e-e9e5-3139-8f0d-fba4c0a2187f",
	Name:        "Main Floor",
	DownloadURL: "",
}

func (suite *AppTestSuite) TestGetOneFloor() {
	floorId := MainFloor.ID
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneFloor,
		URL:     fmt.Sprintf("/floors/%s", floorId),
		URLParams: map[string]string{
			"id": floorId,
		},
	})
	assert.Equal(t, http.StatusOK, rr.Code, "status code")

	var payload *model.Floor
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.True(t, MainFloor.Equal(payload), "incorrect response object")
}

func (suite *AppTestSuite) TestGetOneFloorFail() {
	floorId := "2"
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneFloor,
		URL:     fmt.Sprintf("/floors/%s", floorId),
		URLParams: map[string]string{
			"id": floorId,
		},
	})
	assert.Equal(t, http.StatusNotFound, rr.Code, "status code")
}

func (suite *AppTestSuite) TestGetAllFloors() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetAllFloors,
		URL:     "/floors/",
	})
	assert.Equal(t, http.StatusOK, rr.Code, "status code")

	var payload []*model.Floor
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 1, len(payload), "incorrect response size")
	assert.True(t, MainFloor.Equal(payload[0]), "incorrect response object")
}

// Run after
// https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
func (suite *AppTestSuite) Test_CreateFloors() {
	floorName := "test-floor"
	fileId := "test-file-id"
	mockDrive := new(mockDrive)
	mockDrive.On("UploadFloorPlan", floorName).Return(fileId, nil)
	suite.app.gDrive = mockDrive

	t := suite.T()
	body := new(bytes.Buffer)
	file, err := os.Open("../test-fixtures/test-img.jpg")
	if err != nil {
		t.Fatalf("failed to open file")
	}
	fileContents, _ := ioutil.ReadAll(file)

	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", floorName)
	part, _ := writer.CreateFormFile("image", "test-img.jpg")
	_, _ = part.Write(fileContents)
	_ = writer.Close()

	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    body,
		Handler: suite.app.CreateFloor,
		URL:     "/floors/",
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	})
	if !assert.Equal(t, http.StatusCreated, rr.Code, "status code") {
		return
	}

	var payload *model.Floor
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	driveLink := fmt.Sprintf(`https://drive.google.com/uc?export=download&id=%s`, fileId)
	assert.Equal(t, driveLink, payload.DownloadURL, "wrong drive url")
	assert.Equal(t, floorName, payload.Name, "wrong floor name")
}

func (suite *AppTestSuite) Test_CreateFloorsEmptyName() {
	floorName := "test-floor"
	fileId := "test-file-id"
	mockDrive := new(mockDrive)
	mockDrive.On("UploadFloorPlan", floorName).Return(fileId, nil)
	suite.app.gDrive = mockDrive

	t := suite.T()
	body := new(bytes.Buffer)
	file, err := os.Open("../test-fixtures/test-img.jpg")
	if err != nil {
		t.Fatalf("failed to open file")
	}
	fileContents, _ := ioutil.ReadAll(file)

	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test-img.jpg")
	_, _ = part.Write(fileContents)
	_ = writer.Close()

	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    body,
		Handler: suite.app.CreateFloor,
		URL:     "/floors/",
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	})
	if !assert.Equal(t, http.StatusBadRequest, rr.Code, "status code") {
		return
	}
}

func (suite *AppTestSuite) Test_CreateFloorsEmptyFile() {
	floorName := "test-floor"
	fileId := "test-file-id"
	mockDrive := new(mockDrive)
	mockDrive.On("UploadFloorPlan", floorName).Return(fileId, nil)
	suite.app.gDrive = mockDrive

	t := suite.T()
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", floorName)

	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    body,
		Handler: suite.app.CreateFloor,
		URL:     "/floors/",
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	})
	if !assert.Equal(t, http.StatusBadRequest, rr.Code, "status code") {
		return
	}
}

func (suite *AppTestSuite) Test_CreateFloorsFailToUpload() {
	floorName := "test-floor"
	mockDrive := new(mockDrive)
	mockDrive.On("UploadFloorPlan", floorName).
		Return("", errors.New("some error"))
	suite.app.gDrive = mockDrive

	t := suite.T()
	body := new(bytes.Buffer)
	file, err := os.Open("../test-fixtures/test-img.jpg")
	if err != nil {
		t.Fatalf("failed to open file")
	}
	fileContents, _ := ioutil.ReadAll(file)

	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", floorName)
	part, _ := writer.CreateFormFile("image", "test-img.jpg")
	_, _ = part.Write(fileContents)
	_ = writer.Close()

	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    body,
		Handler: suite.app.CreateFloor,
		URL:     "/floors/",
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	})
	if !assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code") {
		return
	}

}
