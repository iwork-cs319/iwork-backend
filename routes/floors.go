package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (app *App) RegisterFloorRoutes() {
	//app.router.HandleFunc("/floors/", app.CreateFloor).Methods("POST")
	app.router.HandleFunc("/floors/{id}", app.GetOneFloor).Methods("GET")
	app.router.HandleFunc("/floors/", app.GetAllFloors).Methods("GET")
	//app.router.HandleFunc("/floors/{id}", app.UpdateFloor).Methods("PATCH")
	//app.router.HandleFunc("/floors/{id}", app.DeleteFloor).Methods("DELETE")
}

//func (app *App) CreateFloor(w http.ResponseWriter, r *http.Request) {
//	var newFloor model.Floor
//	reqBody, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Printf("App.CreateFloor - error reading request body %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err = json.Unmarshal(reqBody, &newFloor)
//	if err != nil {
//		log.Printf("App.CreateFloor - error unmarshaling request body %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	err = app.store.FloorProvider.CreateFloor(&newFloor)
//	if err != nil {
//		log.Printf("App.CreateFloor - error creating floor %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusCreated)
//	json.NewEncoder(w).Encode(newFloor)
//}

func (app *App) GetOneFloor(w http.ResponseWriter, r *http.Request) {
	floorID := mux.Vars(r)["id"]

	if floorID == "" {
		log.Printf("App.GetOneFloor - empty floor id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	floor, err := app.store.FloorProvider.GetOneFloor(floorID)
	if err != nil {
		log.Printf("App.GetOneFloor - error getting floor from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(floor)
}

func (app *App) GetAllFloors(w http.ResponseWriter, r *http.Request) {
	floors, err := app.store.FloorProvider.GetAllFloors()
	if err != nil {
		log.Printf("App.GetAllFloors - error getting all floors from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(floors)
}

//func (app *App) UpdateFloor(w http.ResponseWriter, r *http.Request) {
//	floorID := mux.Vars(r)["id"]
//
//	if floorID == "" {
//		log.Printf("App.UpdateFloor - empty floor id")
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	var updatedFloor model.Floor
//	reqBody, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Printf("App.UpdateFloor - error reading request body %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	err = json.Unmarshal(reqBody, &updatedFloor)
//	if err != nil {
//		log.Printf("App.UpdateFloor - error unmarshaling request body %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err = app.store.FloorProvider.UpdateFloor(floorID, &updatedFloor)
//	if err != nil {
//		log.Printf("App.UpdateFloor - error getting all floors from provider %v", err)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//}

//func (app *App) DeleteFloor(w http.ResponseWriter, r *http.Request) {
//	floorID := mux.Vars(r)["id"]
//
//	if floorID == "" {
//		log.Printf("App.DeleteFloor - empty floor id")
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err := app.store.FloorProvider.RemoveFloor(floorID)
//	if err != nil {
//		log.Printf("App.DeleteFloor - error getting all floors from provider %v", err)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//}
