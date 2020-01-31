package routes

import (
	"net/http"
)

func (app *App) RegisterBookingRoutes() {
	app.router.HandleFunc("/bookings", app.CreateBooking).Methods("POST")
	app.router.HandleFunc("/bookings/{id}", app.GetOneBooking).Methods("GET")
	app.router.HandleFunc("/bookings/workspaces/{workspace_id}", app.GetBookingsByWorkspaceID).Methods("GET")
	app.router.HandleFunc("/bookings/users/{user_id}", app.GetBookingsByUserID).Methods("GET")
	app.router.HandleFunc("/bookings/start/{start_time}/end/{end_time}", app.GetBookingsByDateRange).Methods("GET")
	app.router.HandleFunc("/bookings", app.GetAllBookings).Methods("GET")
	app.router.HandleFunc("/bookings/{id}", app.UpdateBooking).Methods("PATCH")
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

	err = app.store.BookingProvider.CreateBooking(&newBooking)
	if err != nil {
		log.Printf("App.CreateBooking - error creating booking %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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

	booking, err := app.store.BookingProvider.GetOneBooking(bookingID)
	if err != nil {
		log.Printf("App.GetOneBooking - error getting booking from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

func (app *App) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := app.store.BookingProvider.GetAllBookings()
	if err != nil {
		log.Printf("App.GetAllBookings - error getting all bookings from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bookings)
}

func (app *App) GetBookingsByWorkspaceID(w http.ResponseWriter, r *http.Request) {
	bookings, err := app.store.BookingProvider.GetBookingsByWorkspaceID()
	if err != nil {
		log.Printf("App.GetBookingsByWorkspaceID - error getting bookings by workspaceID from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bookings)
}

func (app *App) GetBookingsByUserID(w http.ResponseWriter, r *http.Request) {
	bookings, err := app.store.BookingProvider.GetBookingsByUserID()
	if err != nil {
		log.Printf("App.GetBookingsByUserID - error getting bookings by userID from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bookings)
}

func (app *App) GetBookingsByDateRange(w http.ResponseWriter, r *http.Request) {
	bookings, err := app.store.BookingProvider.GetBookingsByDateRange()
	if err != nil {
		log.Printf("App.GetBookingsByDateRange - error getting bookings by date range from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bookings)
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
