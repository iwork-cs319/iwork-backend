package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/mail"
	"go-api/model"
	"go-api/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (app *App) RegisterBookingRoutes() {
	app.router.HandleFunc("/bookings", app.CreateBooking).Methods("POST")
	app.router.HandleFunc("/bookings/{id}", app.GetOneBooking).Methods("GET").Queries("expand", "{expand}")
	app.router.HandleFunc("/bookings/{id}", app.GetOneBooking).Methods("GET") // Handles when Query is empty
	app.router.HandleFunc("/bookings/workspaces/{workspace_id}", app.GetBookingsByWorkspaceID).Methods("GET").Queries("expand", "{expand}")
	app.router.HandleFunc("/bookings/workspaces/{workspace_id}", app.GetBookingsByWorkspaceID).Methods("GET") // Handles when Query is empty
	app.router.HandleFunc("/bookings/users/{user_id}", app.GetBookingsByUserID).Methods("GET").Queries("expand", "{expand}")
	app.router.HandleFunc("/bookings/users/{user_id}", app.GetBookingsByUserID).Methods("GET") // Handles when Query is empty
	app.router.HandleFunc("/bookings", app.GetBookingsByDateRange).Methods("GET").Queries("start", "{start:[0-9]+}").Queries("end", "{end:[0-9]+}")
	app.router.HandleFunc("/bookings", app.GetAllBookings).Methods("GET").Queries("expand", "{expand}")
	app.router.HandleFunc("/bookings", app.GetAllBookings).Methods("GET") // Handles when Query is empty
	//app.router.HandleFunc("/bookings/{id}", app.UpdateBooking).Methods("PATCH")
	app.router.HandleFunc("/bookings/{id}", app.RemoveBooking).Methods("DELETE")
}

func (app *App) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var newBooking model.Booking
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("App.CreateBooking - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &newBooking)
	if err != nil {
		log.Printf("App.CreateBooking - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if newBooking.CreatedBy == "" {
		newBooking.CreatedBy = newBooking.UserID
	}
	id, err := app.store.BookingProvider.CreateBooking(&newBooking)
	if err != nil {
		log.Printf("App.CreateBooking - error creating booking %v", err)
		if strings.Contains(err.Error(), "invalid") {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	newBooking.ID = id

	user, err1 := app.store.UserProvider.GetOneUser(newBooking.UserID)
	// todo check for user and send to creator as well
	eBooking, err2 := app.store.BookingProvider.GetOneExpandedBooking(id)
	if err1 == nil && err2 == nil {
		_ = app.email.SendConfirmation(
			mail.Booking,
			&mail.EmailParams{
				Name:          user.Name,
				Email:         user.Email,
				WorkspaceName: eBooking.WorkspaceName,
				FloorName:     eBooking.FloorName,
				Start:         eBooking.StartDate,
				End:           eBooking.EndDate,
			},
		)
	} else {
		log.Printf("Error getting user: %+v; Error getting booking %+v", err1, err2)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBooking)
}

func (app *App) GetOneBooking(w http.ResponseWriter, r *http.Request) {
	bookingID := mux.Vars(r)["id"]
	if bookingID == "" {
		log.Printf("App.GetOneBooking - empty booking id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exp := r.FormValue("expand")
	var expandBool = false
	if exp != "" {
		expand, err := strconv.ParseBool(exp)
		if err != nil {
			log.Printf("App.GetOneBooking - error converting string to boolean from query parameter %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expandBool = expand
	}
	if expandBool == true {
		expandedBooking, err := app.store.BookingProvider.GetOneExpandedBooking(bookingID)
		if err != nil {
			log.Printf("App.GetOneExpandedBooking - error getting expanded booking from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expandedBooking)
	} else {
		booking, err := app.store.BookingProvider.GetOneBooking(bookingID)
		if err != nil {
			log.Printf("App.GetOneBooking - error getting booking from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(booking)
	}
}

func (app *App) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	exp := r.FormValue("expand")
	var expandBool = false
	if exp != "" {
		expand, err := strconv.ParseBool(exp)
		if err != nil {
			log.Printf("App.GetOneBooking - error converting string to boolean from query parameter %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expandBool = expand
	}
	if expandBool == true {
		expandedBookings, err := app.store.BookingProvider.GetAllExpandedBookings()
		if err != nil {
			log.Printf("App.GetAllExpandedBookings - error getting all expanded bookings from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expandedBookings)
	} else {
		bookings, err := app.store.BookingProvider.GetAllBookings()
		if err != nil {
			log.Printf("App.GetAllBookings - error getting all bookings from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(bookings)
	}
}

func (app *App) GetBookingsByWorkspaceID(w http.ResponseWriter, r *http.Request) {
	workspaceID := mux.Vars(r)["workspace_id"]

	if workspaceID == "" {
		log.Printf("App.GetOneBooking - empty booking id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exp := r.FormValue("expand")
	var expandBool = false
	if exp != "" {
		expand, err := strconv.ParseBool(exp)
		if err != nil {
			log.Printf("App.GetOneBooking - error converting string to boolean from query parameter %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expandBool = expand
	}
	if expandBool == true {
		expandedBooking, err := app.store.BookingProvider.GetExpandedBookingsByWorkspaceID(workspaceID)
		if err != nil {
			log.Printf("App.GetExpandedBookingsByWorkspaceID - error getting expanded bookings by workspaceID from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expandedBooking)
	} else {
		bookings, err := app.store.BookingProvider.GetBookingsByWorkspaceID(workspaceID)
		if err != nil {
			log.Printf("App.GetBookingsByWorkspaceID - error getting bookings by workspaceID from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(bookings)
	}
}

func (app *App) GetBookingsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]

	if userID == "" {
		log.Printf("App.GetOneBooking - empty booking id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exp := r.FormValue("expand")
	var expandBool = false
	if exp != "" {
		expand, err := strconv.ParseBool(exp)
		if err != nil {
			log.Printf("App.GetOneBooking - error converting string to boolean from query parameter %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expandBool = expand
	}
	if expandBool == true {
		expandedBooking, err := app.store.BookingProvider.GetExpandedBookingsByUserID(userID)
		if err != nil {
			log.Printf("App.GetExpandedBookingsByUserID - error getting expanded bookings by userID from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expandedBooking)
	} else {
		bookings, err := app.store.BookingProvider.GetBookingsByUserID(userID)
		if err != nil {
			log.Printf("App.GetBookingsByUserID - error getting bookings by userID from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(bookings)
	}
}

func (app *App) GetBookingsByDateRange(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	start := queryParams["start"][0]
	end := queryParams["end"][0]
	startTime, errStart := utils.TimeStampToTime(start) // Unix Timestamp
	endTime, errEnd := utils.TimeStampToTime(end)
	if errStart != nil {
		log.Printf("App.GetBookingsByDateRange - empty start time param: %v", errStart)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errEnd != nil {
		log.Printf("App.GetBookingsByDateRange - empty end time param: %v", errEnd)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exp := r.FormValue("expand")
	var expandBool = false
	if exp != "" {
		expand, err := strconv.ParseBool(exp)
		if err != nil {
			log.Printf("App.GetOneBooking - error converting string to boolean from query parameter %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expandBool = expand
	}
	if expandBool == true {
		expandedBooking, err := app.store.BookingProvider.GetExpandedBookingsByDateRange(startTime, endTime)
		if err != nil {
			log.Printf("App.GetExpandedBookingsByDateRange - error getting expanded bookings by date range from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expandedBooking)
	} else {
		bookings, err := app.store.BookingProvider.GetBookingsByDateRange(startTime, endTime)
		if err != nil {
			log.Printf("App.GetBookingsByDateRange - error getting bookings by date range from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(bookings)
	}
}

func (app *App) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	bookingID := mux.Vars(r)["id"]

	if bookingID == "" {
		log.Printf("App.UpdateBooking - empty booking id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var updatedBooking model.Booking
	reqBody, err := ioutil.ReadAll(r.Body)
	updatedBooking.ID = bookingID
	if err != nil {
		log.Printf("App.UpdateBooking - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &updatedBooking)
	if err != nil {
		log.Printf("App.UpdateBooking - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = app.store.BookingProvider.UpdateBooking(bookingID, &updatedBooking)
	if err != nil {
		log.Printf("App.UpdateBooking - error getting all bookings from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedBooking)
}

func (app *App) RemoveBooking(w http.ResponseWriter, r *http.Request) {
	bookingID := mux.Vars(r)["id"]

	if bookingID == "" {
		log.Printf("App.RemoveBooking - empty booking id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.BookingProvider.RemoveBooking(bookingID)
	if err != nil {
		log.Printf("App.RemoveBooking - error getting all bookings from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
