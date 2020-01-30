package routes

import (
	"net/http"
)

func (app *App) RegisterBookingRoutes() {
	app.router.HandleFunc("/bookings", app.CreateBooking).Methods("POST")
	app.router.HandleFunc("/bookings/{id}", app.GetBooking).Methods("GET")
	app.router.HandleFunc("/bookings", app.GetAllBookings).Methods("GET")
}

func (app *App) CreateBooking(w http.ResponseWriter, r *http.Request) {
	panic("TODO implement me")
}

func (app *App) GetBooking(w http.ResponseWriter, r *http.Request) {
	panic("TODO implement me")
}

func (app *App) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	panic("TODO implement me")
}
