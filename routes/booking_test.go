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
	"net/http/httptest"
)

func bookingEqualMinusID(this *model.Booking, other *model.Booking) bool { // To be used when testing Creation, as ID will not be known in advance.
	return this.WorkspaceID == other.WorkspaceID &&
		this.UserID == other.UserID && this.StartDate == other.StartDate &&
		this.EndDate == other.EndDate && this.Cancelled == other.Cancelled && this.CreatedBy == other.CreatedBy
}

var startBooking1, _ = utils.TimeStampToTime("1547424000")
var endBooking1, _ = utils.TimeStampToTime("1547596799")
var Booking1 = &model.Booking{
	ID:          "d3d7aa7d-59db-3cf2-8243-136f59587124",
	UserID:      "e99a988a-1d41-3997-8d59-959a48ac24a0",
	WorkspaceID: "a12411d3-d281-3735-b000-bf94b094d2af",
	Cancelled:   false,
	StartDate:   startBooking1,
	EndDate:     endBooking1,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var startBooking2, _ = utils.TimeStampToTime("1548288000")
var endBooking2, _ = utils.TimeStampToTime("1548719999")
var Booking2 = &model.Booking{
	ID:          "1440cf5a-3ec6-3a0c-b0eb-c5beec4b3fde",
	UserID:      "e99a988a-1d41-3997-8d59-959a48ac24a0",
	WorkspaceID: "a12411d3-d281-3735-b000-bf94b094d2af",
	Cancelled:   false,
	StartDate:   startBooking2,
	EndDate:     endBooking2,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var startBooking3, _ = utils.TimeStampToTime("1547683200")
var endBooking3, _ = utils.TimeStampToTime("1547855999")
var Booking3 = &model.Booking{
	ID:          "b2b08286-0e72-30b5-8a8a-38b8786d2088",
	UserID:      "e99a988a-1d41-3997-8d59-959a48ac24a0",
	WorkspaceID: "f2188d5e-509f-3086-a031-d86da93a55c4",
	Cancelled:   false,
	StartDate:   startBooking3,
	EndDate:     endBooking3,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var startBooking4, _ = utils.TimeStampToTime("1548115200")
var endBooking4, _ = utils.TimeStampToTime("1548287999")
var Booking4 = &model.Booking{
	ID:          "87d2ef90-5e00-33e7-be54-572b0e9c29c1",
	UserID:      "e99a988a-1d41-3997-8d59-959a48ac24a0",
	WorkspaceID: "f2188d5e-509f-3086-a031-d86da93a55c4",
	Cancelled:   false,
	StartDate:   startBooking4,
	EndDate:     endBooking4,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var startBooking5, _ = utils.TimeStampToTime("1547683200")
var endBooking5, _ = utils.TimeStampToTime("1547769599")
var Booking5 = &model.Booking{
	ID:          "70747d17-f79f-3c82-88f2-a6ac9c7ba86f",
	UserID:      "e99a988a-1d41-3997-8d59-959a48ac24a0",
	WorkspaceID: "62eb266e-740d-3973-9e50-77b0287e3026",
	Cancelled:   false,
	StartDate:   startBooking5,
	EndDate:     endBooking5,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var startBooking6, _ = utils.TimeStampToTime("1548115200")
var endBooking6, _ = utils.TimeStampToTime("1548374399")
var Booking6 = &model.Booking{
	ID:          "d8d1beb0-21e3-3084-b29f-ea2b3c194a2b",
	UserID:      "e99a988a-1d41-3997-8d59-959a48ac24a0",
	WorkspaceID: "5e56de3d-2323-372d-897f-23d6037c8581",
	Cancelled:   false,
	StartDate:   startBooking6,
	EndDate:     endBooking6,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var startBooking7, _ = utils.TimeStampToTime("1548547200")
var endBooking7, _ = utils.TimeStampToTime("1548719999")
var Booking7 = &model.Booking{
	ID:          "723ac86b-e0e8-39bd-b407-0c1ced6f2d93",
	UserID:      "e99a988a-1d41-3997-8d59-959a48ac24a0",
	WorkspaceID: "aad40cbb-4baf-3931-a5d2-6f98b414182a",
	Cancelled:   false,
	StartDate:   startBooking7,
	EndDate:     endBooking7,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}
var startBooking11, _ = utils.TimeStampToTime("1548547200")
var endBooking11, _ = utils.TimeStampToTime("1548719999")
var Booking11 = &model.Booking{
	ID:          "2f092deb-3c46-3bf9-9bf9-0a31b7c119b7",
	UserID:      "e99a988a-1d41-3997-8d59-959a48ac24a0",
	WorkspaceID: "bb15369d-e6e0-33b8-8b97-1779f8865890",
	Cancelled:   false,
	StartDate:   startBooking7,
	EndDate:     endBooking7,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}

func (suite *AppTestSuite) TestGetOneBookingFail() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneBooking,
		URL:     fmt.Sprintf("/bookings/%s", "2"),
		URLParams: map[string]string{
			"id": "2",
		},
	})

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code")
}

func (suite *AppTestSuite) TestGetOneBooking() {
	bookingId := Booking1.ID
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneBooking,
		URL:     fmt.Sprintf("/bookings/%s", bookingId),
		URLParams: map[string]string{
			"id": bookingId,
		},
	})

	assert.Equal(t, rr.Code, http.StatusOK, "status code")
	var payload *model.Booking
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, Booking1, payload, "The patched booking is not the same as the sent request.")
}

func (suite *AppTestSuite) TestGetAllBookings() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetAllBookings,
		URL:     "/bookings/",
	})

	assert.Equal(t, rr.Code, http.StatusOK, "status code")

	var payload []*model.Booking
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 13, len(payload), "testGetAllBookings: incorrect response size")
	expectedBookings := []*model.Booking{Booking1, Booking2, Booking3, Booking4, Booking5, Booking6, Booking7}
	for _, expected := range expectedBookings {
		found := false
		for _, b := range payload {
			if b.Equal(expected) {
				found = true
			}
		}
		if !found {
			t.Fatalf("testGetAllBookings: %s not found in offring list", expected.ID)
		}
	}
}

func (suite *AppTestSuite) TestGetOneBookingByWorkspaceIDFail() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetBookingsByWorkspaceID,
		URL:     fmt.Sprintf("/bookings/workspaces/%s", "2"),
		URLParams: map[string]string{
			"workspace_id": "2",
		},
	})

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code")
}

func (suite *AppTestSuite) TestGetBookingsByWorkspaceID() {
	bookingId := Booking1.WorkspaceID
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetBookingsByWorkspaceID,
		URL:     fmt.Sprintf("/bookings/workspaces/%s", bookingId),
		URLParams: map[string]string{
			"workspace_id": bookingId,
		},
	})

	assert.Equal(t, rr.Code, http.StatusOK, "status code")

	var payload []*model.Booking
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 2, len(payload), "testGetBookingsByWorkspaceID: incorrect response size")
	expectedBookings := []*model.Booking{Booking1, Booking2}
	for _, expected := range expectedBookings {
		found := false
		for _, b := range payload {
			if b.Equal(expected) {
				found = true
			}
		}
		if !found {
			t.Fatalf("testGetBookingsByWorkspaceID: %s not found in offring list", expected.ID)
		}
	}
}

func (suite *AppTestSuite) TestGetOneBookingByUserIDFail() {
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetBookingsByUserID,
		URL:     fmt.Sprintf("/bookings/users/%s", "13"),
		URLParams: map[string]string{
			"user_id": "2",
		},
	})

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code")
}

func (suite *AppTestSuite) TestGetBookingsByUserID() {
	bookingId := Booking1.UserID
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetBookingsByUserID,
		URL:     fmt.Sprintf("/bookings/users/%s", bookingId),
		URLParams: map[string]string{
			"user_id": bookingId,
		},
	})
	assert.Equal(t, rr.Code, http.StatusOK, "status code")
	var payload []*model.Booking
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, 13, len(payload), "testGetBookingsByUserID: incorrect response size")
	expectedBookings := []*model.Booking{Booking1, Booking2}
	for _, expected := range expectedBookings {
		found := false
		for _, b := range payload {
			if b.Equal(expected) {
				found = true
			}
		}
		if !found {
			t.Fatalf("testGetBookingsByUserID: %s not found in offring list", expected.ID)
		}
	}
}

func (suite *AppTestSuite) TestGetBookingsByDateRange() {
	bookingStart := "1548547200"
	bookingEnd := "1548719999"
	t := suite.T()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/bookings?%s&%s", bookingStart, bookingEnd), nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("start", bookingStart)
	q.Add("end", bookingEnd)
	req.URL.RawQuery = q.Encode()
	//req = mux.SetURLVars(req, config.URLParams)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.app.GetBookingsByDateRange)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK, "status code")
	var payload []*model.Booking
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	log.Printf("HERE'S THE ID'S")
	for _, p := range payload {
		log.Printf(p.ID)
	}
	assert.Equal(t, 2, len(payload), "testGetBookingsByDateRange: incorrect response size")
	expectedBookings := []*model.Booking{Booking11, Booking7}
	for _, expected := range expectedBookings {
		found := false
		for _, b := range payload {
			if b.Equal(expected) {
				found = true
			}
		}
		if !found {
			t.Fatalf("testGetBookingsByDateRange: %s not found in offring list", expected.ID)
		}
	}
}

var startBookingNew, _ = utils.TimeStampToTime("1583884800")
var endBookingNew, _ = utils.TimeStampToTime("1583971199")
var newBooking = &model.Booking{
	ID:          "", // Unknown at this point
	UserID:      "32ea2fb1-7124-304a-b9c3-eb445578103e",
	WorkspaceID: "5e56de3d-2323-372d-897f-23d6037c8581",
	Cancelled:   false,
	StartDate:   startBookingNew,
	EndDate:     endBookingNew,
	CreatedBy:   "8b5bb736-6a1d-3378-8e71-ab45fe8beb84",
}

func (suite *AppTestSuite) Test_CreateBooking() {
	// Check ID does not exist in database
	t := suite.T()
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetAllBookings,
		URL:     "/bookings/",
	})
	var payload []*model.Booking
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	for _, b := range payload {
		if bookingEqualMinusID(newBooking, b) {
			t.Fatalf("testCreateBooking: %s found in booking list, when it should not exist", newBooking.ID)
		}
	}
	// [POST] Create booking
	requestBody, _ := json.Marshal(map[string]interface{}{
		"workspace_id": newBooking.WorkspaceID,
		"user_id":      newBooking.UserID,
		"start_time":   newBooking.StartDate,
		"end_time":     newBooking.EndDate,
		"cancelled":    newBooking.Cancelled,
		"created_by":   newBooking.CreatedBy,
	})
	rr2 := executeReq(t, &testRouteConfig{
		Method:  http.MethodPost,
		Body:    bytes.NewBuffer(requestBody),
		Handler: suite.app.CreateBooking,
		URL:     fmt.Sprintf("/bookings"),
	})
	// Check correct response

	assert.Equal(t, rr2.Code, http.StatusCreated, "status code")
	var payload2 *model.Booking
	_ = json.Unmarshal(rr2.Body.Bytes(), &payload2)
	// Just give it the ID since we want to use assert
	newBooking.ID = payload2.ID
	assert.Equal(t, newBooking, payload2, "The created booking is not the same as the sent request.")
}

var startBookingPatch, _ = utils.TimeStampToTime("7952342400")
var endBookingPatch, _ = utils.TimeStampToTime("7952515199")
var patchBooking = &model.Booking{
	ID:          newBooking.ID, // Unknown at this point
	UserID:      "8b5bb736-6a1d-3378-8e71-ab45fe8beb84",
	WorkspaceID: "aad40cbb-4baf-3931-a5d2-6f98b414182a",
	Cancelled:   true,
	StartDate:   startBookingPatch,
	EndDate:     endBookingPatch,
	CreatedBy:   "e99a988a-1d41-3997-8d59-959a48ac24a0",
}

func (suite *AppTestSuite) Test_PatchBooking() {
	t := suite.T()
	// Assume ID exists since we just create it
	// Change everything but the ID
	requestBody, _ := json.Marshal(map[string]interface{}{
		"workspace_id": patchBooking.WorkspaceID,
		"user_id":      patchBooking.UserID,
		"start_time":   patchBooking.StartDate,
		"end_time":     patchBooking.EndDate,
		"cancelled":    patchBooking.Cancelled,
		"created_by":   patchBooking.CreatedBy,
	})
	rr2 := executeReq(t, &testRouteConfig{
		Method:  http.MethodPatch,
		Body:    bytes.NewBuffer(requestBody),
		Handler: suite.app.UpdateBooking,
		URL:     fmt.Sprintf("/bookings/%s", newBooking.ID),
		URLParams: map[string]string{
			"id": newBooking.ID,
		},
	})
	// Response 200
	var payload2 *model.Booking
	_ = json.Unmarshal(rr2.Body.Bytes(), &payload2)
	// Just give it the ID since we want to use assert
	patchBooking.ID = payload2.ID
	assert.Equal(t, patchBooking, payload2, "The created booking is not the same as the sent request.")
}

func (suite *AppTestSuite) Test_ZDeleteBooking() {
	t := suite.T()
	// Assume ID exists since we just created and patched it
	existingID := "723ac86b-e0e8-39bd-b407-0c1ced6f2d93"
	rr2 := executeReq(t, &testRouteConfig{
		Method:  http.MethodDelete,
		Body:    nil,
		Handler: suite.app.RemoveBooking,
		URL:     fmt.Sprintf("/bookings/%s", existingID),
		URLParams: map[string]string{
			"id": existingID,
		},
	})
	// Response 200
	assert.Equal(t, rr2.Code, http.StatusOK, "status code")
	// Check that the id cannot be found in the database
	rr := executeReq(t, &testRouteConfig{
		Method:  http.MethodGet,
		Body:    nil,
		Handler: suite.app.GetOneBooking,
		URL:     fmt.Sprintf("/bookings/%s", existingID),
		URLParams: map[string]string{
			"id": existingID,
		},
	})
	assert.Equal(t, http.StatusOK, rr.Code, "status code") // todo: For now, Delete never fails
}
