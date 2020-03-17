package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/model"
	"go-api/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func (app *App) RegisterOfferingRoutes() {
	app.router.HandleFunc("/offerings", app.CreateOffering).Methods("POST")
	app.router.HandleFunc("/offerings/{id}", app.GetOneOffering).Methods("GET").Queries("expand", "{expand}")
	app.router.HandleFunc("/offerings/{id}", app.GetOneOffering).Methods("GET")
	app.router.HandleFunc("/offerings/workspaces/{workspace_id}", app.GetOfferingsByWorkspaceID).Methods("GET").Queries("expand", "{expand}")
	app.router.HandleFunc("/offerings/workspaces/{workspace_id}", app.GetOfferingsByWorkspaceID).Methods("GET")
	app.router.HandleFunc("/offerings/users/{user_id}", app.GetOfferingsByUserID).Methods("GET").Queries("expand", "{expand}")
	app.router.HandleFunc("/offerings/users/{user_id}", app.GetOfferingsByUserID).Methods("GET")
	app.router.HandleFunc("/offerings", app.GetOfferingsByDateRange).Methods("GET").Queries("start", "{start:[0-9]+}").Queries("end", "{end:[0-9]+}").Queries("expand", "{expand}")
	app.router.HandleFunc("/offerings", app.GetOfferingsByDateRange).Methods("GET").Queries("start", "{start:[0-9]+}").Queries("end", "{end:[0-9]+}")
	app.router.HandleFunc("/offerings", app.GetAllOfferings).Methods("GET").Queries("expand", "{expand}")
	app.router.HandleFunc("/offerings", app.GetAllOfferings).Methods("GET")
	//app.router.HandleFunc("/offerings/{id}", app.UpdateOffering).Methods("PATCH")
	app.router.HandleFunc("/offerings/{id}", app.RemoveOffering).Methods("DELETE")
}

func (app *App) CreateOffering(w http.ResponseWriter, r *http.Request) {
	var newOffering model.Offering
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("App.CreateOffering - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &newOffering)
	if err != nil {
		log.Printf("App.CreateOffering - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if newOffering.CreatedBy == "" {
		newOffering.CreatedBy = newOffering.UserID
	}
	boolean, err := app.store.AssigneeProvider.IsFullyAssigned(newOffering.WorkspaceID)
	if err != nil {
		log.Printf("App.CreateOffering - error checking assignment with workspace ID %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !boolean { // If it isn't fully assigned
		log.Printf("App.CreateOffering - error cannot create assignment on non-assigned workspace %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = app.store.OfferingProvider.GetOfferingsByWorkspaceIDAndDateRange(newOffering.WorkspaceID, newOffering.StartDate, newOffering.EndDate)
	if err != nil {
		log.Printf("App.CreateOffering - error cannot create offering, it already exists! %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := app.store.OfferingProvider.CreateOffering(&newOffering)
	if err != nil {
		log.Printf("App.CreateOffering - error creating offering %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newOffering.ID = id

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOffering)
}

func (app *App) GetOneOffering(w http.ResponseWriter, r *http.Request) {
	offeringID := mux.Vars(r)["id"]

	if offeringID == "" {
		log.Printf("App.GetOneOffering - empty offering id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exp := r.FormValue("expand")
	var expandBool = false
	if exp != "" {
		expand, err := strconv.ParseBool(exp)
		if err != nil {
			log.Printf("App.GetOneOffering - error converting string to boolean from query parameter %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expandBool = expand
	}
	if expandBool == true {
		expandedOffering, err := app.store.OfferingProvider.GetOneExpandedOffering(offeringID)
		if err != nil {
			log.Printf("App.GetOneExpandedOffering - error getting expanded booking from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expandedOffering)
	} else {
		offering, err := app.store.OfferingProvider.GetOneOffering(offeringID)
		if err != nil {
			log.Printf("App.GetOneOffering - error getting offering from provider %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(offering)
	}
}

func (app *App) GetAllOfferings(w http.ResponseWriter, r *http.Request) {
	exp := r.FormValue("expand")
	var expandBool = false
	if exp != "" {
		expand, err := strconv.ParseBool(exp)
		if err != nil {
			log.Printf("App.GetAllOfferings - error converting string to boolean from query parameter %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expandBool = expand
	}
	if expandBool == true {
		expandedOfferings, err := app.store.OfferingProvider.GetAllExpandedOfferings()
		if err != nil {
			log.Printf("App.GetAllExpandedOfferings - error getting all expanded offerings from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(expandedOfferings)
	} else {
		offerings, err := app.store.OfferingProvider.GetAllOfferings()
		if err != nil {
			log.Printf("App.GetAllOfferings - error getting all offerings from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(offerings)
	}
}

func (app *App) GetOfferingsByWorkspaceID(w http.ResponseWriter, r *http.Request) {
	workspaceID := mux.Vars(r)["workspace_id"]
	if workspaceID == "" {
		log.Printf("App.GetOneOffering - empty offering id")
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
		expandedOfferings, err := app.store.OfferingProvider.GetExpandedOfferingsByWorkspaceID(workspaceID)
		if err != nil {
			log.Printf("App.GetExpandedOfferingsByWorkspaceID - error getting expanded offerings by workspaceID from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(expandedOfferings)
	} else {
		offerings, err := app.store.OfferingProvider.GetOfferingsByWorkspaceID(workspaceID)
		if err != nil {
			log.Printf("App.GetOfferingsByWorkspaceID - error getting offerings by workspaceID from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(offerings)
	}
}

func (app *App) GetOfferingsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]

	if userID == "" {
		log.Printf("App.GetOneOffering - empty offering id")
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
		expandedOfferings, err := app.store.OfferingProvider.GetExpandedOfferingsByUserID(userID)
		if err != nil {
			log.Printf("App.GetExpandedOfferingsByUserID - error getting expanded offerings by userID from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(expandedOfferings)
	} else {
		offerings, err := app.store.OfferingProvider.GetOfferingsByUserID(userID)
		if err != nil {
			log.Printf("App.GetOfferingsByUserID - error getting offerings by userID from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(offerings)
	}
}

func (app *App) GetOfferingsByDateRange(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	start := queryParams["start"][0]
	end := queryParams["end"][0]
	startTime, errStart := utils.TimeStampToTime(start) // Unix Timestamp
	endTime, errEnd := utils.TimeStampToTime(end)
	if errStart != nil {
		log.Printf("App.GetOfferingsByDateRange - empty start time param: %v", errStart)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errEnd != nil {
		log.Printf("App.GetOfferingsByDateRange - empty end time param: %v", errEnd)
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
		expandedOfferings, err := app.store.OfferingProvider.GetExpandedOfferingsByDateRange(startTime, endTime)
		if err != nil {
			log.Printf("App.GetExpandedOfferingsByDateRange - error getting expanded offerings by date range from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(expandedOfferings)
	} else {
		offerings, err := app.store.OfferingProvider.GetOfferingsByDateRange(startTime, endTime)
		if err != nil {
			log.Printf("App.GetOfferingsByDateRange - error getting offerings by date range from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(offerings)
	}
}

func (app *App) UpdateOffering(w http.ResponseWriter, r *http.Request) {
	offeringID := mux.Vars(r)["id"]

	if offeringID == "" {
		log.Printf("App.UpdateOffering - empty offering id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var updatedOffering model.Offering
	reqBody, err := ioutil.ReadAll(r.Body)
	updatedOffering.ID = offeringID
	if err != nil {
		log.Printf("App.UpdateOffering - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &updatedOffering)
	if err != nil {
		log.Printf("App.UpdateOffering - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = app.store.OfferingProvider.UpdateOffering(offeringID, &updatedOffering)
	if err != nil {
		log.Printf("App.UpdateOffering - error getting all offerings from provider %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedOffering)
}

func (app *App) RemoveOffering(w http.ResponseWriter, r *http.Request) {
	offeringID := mux.Vars(r)["id"]

	if offeringID == "" {
		log.Printf("App.RemoveOffering - empty offering id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	offering, err := app.store.OfferingProvider.GetOneOffering(offeringID)
	if err != nil {
		log.Printf("App.RemoveOffering - offering does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	boolean, err := app.store.AssigneeProvider.IsFullyAssigned(offering.WorkspaceID)
	if !boolean {
		log.Printf("App.RemoveOffering - offering not for assigned seating, cannot be cancelled.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = app.store.OfferingProvider.RemoveOffering(offeringID)
	if err != nil {
		log.Printf("App.RemoveOffering - error getting all offerings from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
