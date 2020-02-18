package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/model"
	"go-api/utils"
	"io/ioutil"
	"log"
	"net/http"
)

func (app *App) RegisterWorkspaceRoutes() {
	app.router.HandleFunc("/workspaces/available", app.GetAvailability).Methods("GET").Queries("start", "{start:[0-9]+}").Queries("end", "{end:[0-9]+}")
	app.router.HandleFunc("/workspaces", app.CreateWorkspace).Methods("POST")
	app.router.HandleFunc("/workspaces/{id}", app.GetOneWorkspace).Methods("GET")
	app.router.HandleFunc("/workspaces", app.GetAllWorkspaces).Methods("GET")
	app.router.HandleFunc("/workspaces/{id}", app.UpdateWorkspace).Methods("PATCH")
	app.router.HandleFunc("/workspaces/{id}", app.DeleteWorkspace).Methods("DELETE")
}

func (app *App) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var newWorkspace model.Workspace
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("App.CreateWorkspace - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &newWorkspace)
	if err != nil {
		log.Printf("App.CreateWorkspace - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := app.store.WorkspaceProvider.CreateWorkspace(&newWorkspace)
	if err != nil {
		log.Printf("App.CreateWorkspace - error creating workspace %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newWorkspace.ID = id

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newWorkspace)
}

func (app *App) GetOneWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceID := mux.Vars(r)["id"]

	if workspaceID == "" {
		log.Printf("App.GetOneWorkspace - empty workspace id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	workspace, err := app.store.WorkspaceProvider.GetOneWorkspace(workspaceID)
	if err != nil {
		log.Printf("App.GetOneWorkspace - error getting workspace from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspace)
}

func (app *App) GetAllWorkspaces(w http.ResponseWriter, r *http.Request) {
	workspaces, err := app.store.WorkspaceProvider.GetAllWorkspaces()
	if err != nil {
		log.Printf("App.GetAllWorkspaces - error getting all workspaces from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(workspaces)
}

func (app *App) UpdateWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceID := mux.Vars(r)["id"]

	if workspaceID == "" {
		log.Printf("App.UpdateWorkspace - empty workspace id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var updatedWorkspace model.Workspace
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("App.UpdateWorkspace - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &updatedWorkspace)
	if err != nil {
		log.Printf("App.UpdateWorkspace - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = app.store.WorkspaceProvider.UpdateWorkspace(workspaceID, &updatedWorkspace)
	if err != nil {
		log.Printf("App.UpdateWorkspace - error getting all workspaces from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (app *App) DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceID := mux.Vars(r)["id"]

	if workspaceID == "" {
		log.Printf("App.DeleteWorkspace - empty workspace id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.WorkspaceProvider.RemoveWorkspace(workspaceID)
	if err != nil {
		log.Printf("App.DeleteWorkspace - error getting all workspaces from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) GetAvailability(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	start := queryParams["start"][0]
	end := queryParams["end"][0]
	startTime, errStart := utils.TimeStampToTime(start) // Unix Timestamp
	endTime, errEnd := utils.TimeStampToTime(end)
	if errStart != nil {
		log.Printf("App.GetAvailability - error getting offerings by date range from provider %v", errStart)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errEnd != nil {
		log.Printf("App.GetAvailability - error getting offerings by date range from provider %v", errEnd)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	workspaceIds, err := app.store.WorkspaceProvider.FindAvailability("e1cb788a-2e37-4950-a9a5-db8612d4cc80", startTime, endTime)
	if err != nil {
		log.Printf("App.GetOfferingsByDateRange - error getting ids from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(workspaceIds)
}
