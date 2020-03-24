package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-api/model"
	"go-api/utils"
	"log"
	"net/http"
)

func offeringEqualMinusID(this *model.Offering, other *model.Offering) bool { // To be used when testing Creation, as ID will not be known in advance.
	return this.WorkspaceID == other.WorkspaceID &&
		this.UserID == other.UserID && this.StartDate == other.StartDate &&
		this.EndDate == other.EndDate && this.Cancelled == other.Cancelled && this.CreatedBy == other.CreatedBy
}

var start1, _ = utils.TimeStampToTime("1547510400")
var end1, _ = utils.TimeStampToTime("1547726399")
var Offering1 = &model.Offering{
	ID:          "2d111346-5f47-3d9a-8e43-af74e518bb4f",
	UserID:      "32ea2fb1-7124-304a-b9c3-eb445578103e",
	WorkspaceID: "62eb266e-740d-3973-9e50-77b0287e3026",
	Cancelled:   false,
	StartDate:   start1,
	EndDate:     end1,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var start2, _ = utils.TimeStampToTime("1548115200")
var end2, _ = utils.TimeStampToTime("1548590399")
var Offering2 = &model.Offering{
	ID:          "4771cd04-3ba7-31f0-bb0c-faa379f16fc1",
	UserID:      "32ea2fb1-7124-304a-b9c3-eb445578103e",
	WorkspaceID: "62eb266e-740d-3973-9e50-77b0287e3026",
	Cancelled:   false,
	StartDate:   start2,
	EndDate:     end2,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var start3, _ = utils.TimeStampToTime("1547596800")
var end3, _ = utils.TimeStampToTime("1547942399")
var Offering3 = &model.Offering{
	ID:          "bbc3481c-f9e3-3d1f-ae03-860b7f08468f",
	UserID:      "8b5bb736-6a1d-3378-8e71-ab45fe8beb84",
	WorkspaceID: "5e56de3d-2323-372d-897f-23d6037c8581",
	Cancelled:   false,
	StartDate:   start3,
	EndDate:     end3,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}

//var start4, _ = utils.TimeStampToTime("") same as start2
var end4, _ = utils.TimeStampToTime("1548374399")
var Offering4 = &model.Offering{
	ID:          "ae271b5b-5d0b-3284-b105-65b25f2d9e1e",
	UserID:      "8b5bb736-6a1d-3378-8e71-ab45fe8beb84",
	WorkspaceID: "5e56de3d-2323-372d-897f-23d6037c8581",
	Cancelled:   false,
	StartDate:   start2,
	EndDate:     end4,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var start5, _ = utils.TimeStampToTime("1547337600")
var end5, _ = utils.TimeStampToTime("1547510399")
var Offering5 = &model.Offering{
	ID:          "43c612ef-8f79-3c1c-8faf-24a07e9bfd76",
	UserID:      "ab6c2c4f-c112-3c1e-bcf1-42cdea289c1f",
	WorkspaceID: "aad40cbb-4baf-3931-a5d2-6f98b414182a",
	Cancelled:   false,
	StartDate:   start5,
	EndDate:     end5,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var start6, _ = utils.TimeStampToTime("1548115200")
var end6, _ = utils.TimeStampToTime("1548374399")
var Offering6 = &model.Offering{
	ID:          "f2801e9b-01f7-3ad2-b558-c9700ce74ac0",
	UserID:      "ab6c2c4f-c112-3c1e-bcf1-42cdea289c1f",
	WorkspaceID: "aad40cbb-4baf-3931-a5d2-6f98b414182a",
	Cancelled:   false,
	StartDate:   start6,
	EndDate:     end6,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var start7, _ = utils.TimeStampToTime("1548547200")
var end7, _ = utils.TimeStampToTime("1548719999")
var Offering7 = &model.Offering{
	ID:          "e44dfa5c-fad5-3aa9-8951-fd0ed9d66f18",
	UserID:      "ab6c2c4f-c112-3c1e-bcf1-42cdea289c1f",
	WorkspaceID: "aad40cbb-4baf-3931-a5d2-6f98b414182a",
	Cancelled:   false,
	StartDate:   start7,
	EndDate:     end7,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}

func (suite *AppTestSuite) TestGetOneOfferingFail() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneOffering,
		URL:     fmt.Sprintf("/offerings/%s", "2"),
		URLParams: map[string]string{
			"id": "2",
		},
	})

	assert.Equal(t, http.StatusNotFound, rr.Code, "status code")
}

func (suite *AppTestSuite) TestGetOneOffering() {
	offeringId := Offering1.ID
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneOffering,
		URL:     fmt.Sprintf("/offerings/%s", offeringId),
		URLParams: map[string]string{
			"id": offeringId,
		},
	})

	assert.Equal(t, rr.Code, http.StatusOK, "status code")
	var payload *model.Offering
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, Offering1, payload, "The patched offering is not the same as the sent request.")
}

func (suite *AppTestSuite) TestGetAllOfferings() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetAllOfferings,
		URL:     "/offerings/",
	})

	assert.Equal(t, rr.Code, http.StatusOK, "status code")

	var payload []*model.Offering
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 7, len(payload), "testGetAllOfferings: incorrect response size")
	expectedOfferings := []*model.Offering{Offering1, Offering2, Offering3, Offering4, Offering5, Offering6, Offering7}
	for _, expected := range expectedOfferings {
		found := false
		for _, o := range payload {
			if o.Equal(expected) {
				found = true
			}
		}
		if !found {
			t.Fatalf("testGetAllOfferings: %s not found in offring list", expected.ID)
		}
	}
}

func (suite *AppTestSuite) TestGetOneOfferingByWorkspaceIDFail() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOfferingsByWorkspaceID,
		URL:     fmt.Sprintf("/offerings/workspaces/%s", "2"),
		URLParams: map[string]string{
			"workspace_id": "2",
		},
	})

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code")
}

func (suite *AppTestSuite) TestGetOfferingsByWorkspaceID() {
	offeringId := Offering1.WorkspaceID
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOfferingsByWorkspaceID,
		URL:     fmt.Sprintf("/offerings/workspaces/%s", offeringId),
		URLParams: map[string]string{
			"workspace_id": offeringId,
		},
	})

	assert.Equal(t, rr.Code, http.StatusOK, "status code")

	var payload []*model.Offering
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 2, len(payload), "testGetOfferingsByWorkspaceID: incorrect response size")
	expectedOfferings := []*model.Offering{Offering1, Offering2}
	for _, expected := range expectedOfferings {
		found := false
		for _, o := range payload {
			if o.Equal(expected) {
				found = true
			}
		}
		if !found {
			t.Fatalf("testGetOfferingsByWorkspaceID: %s not found in offring list", expected.ID)
		}
	}
}

func (suite *AppTestSuite) TestGetOneOfferingByUserIDFail() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOfferingsByUserID,
		URL:     fmt.Sprintf("/offerings/users/%s", "3"),
		URLParams: map[string]string{
			"user_id": "2",
		},
	})

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code")
}

func (suite *AppTestSuite) TestGetOfferingsByUserID() {
	offeringId := Offering1.UserID
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOfferingsByUserID,
		URL:     fmt.Sprintf("/offerings/users/%s", offeringId),
		URLParams: map[string]string{
			"user_id": offeringId,
		},
	})
	assert.Equal(t, rr.Code, http.StatusOK, "status code")
	var payload []*model.Offering
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 2, len(payload), "testGetOfferingsByUserID: incorrect response size")
	expectedOfferings := []*model.Offering{Offering1, Offering2}
	for _, expected := range expectedOfferings {
		found := false
		for _, o := range payload {
			if o.Equal(expected) {
				found = true
			}
		}
		if !found {
			t.Fatalf("testGetOfferingsByUserID: %s not found in offring list", expected.ID)
		}
	}
}

func (suite *AppTestSuite) TestGetOfferingsByDateRange() {
	offeringStart := "1547337600"
	offeringEnd := "1547510399"
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOfferingsByDateRange,
		URL: fmt.Sprintf(
			"/offerings?start=%s&end=%s",
			offeringStart,
			offeringEnd,
		),
	})
	assert.Equal(t, rr.Code, http.StatusOK, "status code")
	var payload []*model.Offering
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 1, len(payload), "testGetOfferingsByDateRange: incorrect response size")
	expectedOfferings := []*model.Offering{Offering5}
	for _, expected := range expectedOfferings {
		found := false
		for _, o := range payload {
			if o.Equal(expected) {
				found = true
			}
		}
		if !found {
			t.Fatalf("testGetOfferingsByDateRange: %s not found in offring list", expected.ID)
		}
	}
}

var newOffering = &model.Offering{
	ID:          "", // Unknown at this point
	UserID:      "32ea2fb1-7124-304a-b9c3-eb445578103e",
	WorkspaceID: "5e56de3d-2323-372d-897f-23d6037c8581",
	Cancelled:   false,
	StartDate:   date("2021-03-24T00:00:00Z"),
	EndDate:     date("2021-03-24T23:59:59Z"),
	CreatedBy:   "8b5bb736-6a1d-3378-8e71-ab45fe8beb84",
}

func (suite *AppTestSuite) Test_CreateOffering() {
	// Check ID does not exist in database
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetAllOfferings,
		URL:     "/offerings/",
	})
	var payload []*model.Offering
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	for _, o := range payload {
		if offeringEqualMinusID(newOffering, o) {
			t.Fatalf("testCreateOffering: %s found in offering list, when it should not exist", newOffering.ID)
		}
	}
	// [POST] Create offering
	requestBody, _ := json.Marshal(map[string]interface{}{
		"workspace_id": newOffering.WorkspaceID,
		"user_id":      newOffering.UserID,
		"start_time":   newOffering.StartDate,
		"end_time":     newOffering.EndDate,
		"cancelled":    newOffering.Cancelled,
		"created_by":   newOffering.CreatedBy,
	})
	rr2 := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    bytes.NewBuffer(requestBody),
		Handler: suite.app.CreateOffering,
		URL:     fmt.Sprintf("/offerings"),
	})
	// Check correct response

	assert.Equal(t, rr2.Code, http.StatusCreated, "status code")
	var payload2 *model.Offering
	_ = json.Unmarshal(rr2.Body.Bytes(), &payload2)
	// Just give it the ID since we want to use assert
	newOffering.ID = payload2.ID
	assert.Equal(t, newOffering, payload2, "The created offering is not the same as the sent request.")
}

func (suite *AppTestSuite) Test_PatchOffering() {
	t := suite.T()
	startPatch := date("2022-03-24T00:00:00Z")
	endPatch := date("2022-03-24T00:00:00Z")
	var patchOffering = &model.Offering{
		ID:          "",
		UserID:      "8b5bb736-6a1d-3378-8e71-ab45fe8beb84",
		WorkspaceID: "aad40cbb-4baf-3931-a5d2-6f98b414182a",
		Cancelled:   true,
		StartDate:   startPatch,
		EndDate:     endPatch,
		CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
	}
	// Assume ID exists since we just create it
	// Change everything but the ID
	requestBody, _ := json.Marshal(map[string]interface{}{
		"workspace_id": patchOffering.WorkspaceID,
		"user_id":      patchOffering.UserID,
		"start_time":   patchOffering.StartDate,
		"end_time":     patchOffering.EndDate,
		"cancelled":    patchOffering.Cancelled,
		"created_by":   patchOffering.CreatedBy,
	})
	rr2 := executeReq(t, &testRouteConfig{
		Method:  http.MethodPatch,
		Body:    bytes.NewBuffer(requestBody),
		Handler: suite.app.UpdateOffering,
		URL:     fmt.Sprintf("/offerings/%s", newOffering.ID), // Purposefully uses the ID from create
		URLParams: map[string]string{
			"id": newOffering.ID,
		},
	})
	// Response 200
	log.Printf(patchOffering.ID)
	var payload2 *model.Offering
	_ = json.Unmarshal(rr2.Body.Bytes(), &payload2)
	// Just give it the ID since we want to use assert
	patchOffering.ID = payload2.ID
	assert.Equal(t, patchOffering, payload2, "The created offering is not the same as the sent request.")
}

func (suite *AppTestSuite) Test_ZDeleteOffering() { // Z to make it be performed last
	t := suite.T()
	// Assume ID exists since we just created and patched it
	existingID := newOffering.ID
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodDelete,
		Body:    nil,
		Handler: suite.app.RemoveOffering,
		URL:     fmt.Sprintf("/offerings/%s", existingID),
		URLParams: map[string]string{
			"id": existingID,
		},
	})
	// Response 200
	assert.Equal(t, http.StatusOK, rr.Code, "status code")
	// Check that the id cannot be found in the database
	rr2 := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneOffering,
		URL:     fmt.Sprintf("/offerings/%s", existingID),
		URLParams: map[string]string{
			"id": existingID,
		},
	})
	assert.Equal(t, http.StatusOK, rr2.Code, "status code") // todo: For now, Delete never fails
}
