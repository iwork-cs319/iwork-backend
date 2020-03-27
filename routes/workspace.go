package routes

import (
	"encoding/csv"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/model"
	"go-api/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (app *App) RegisterWorkspaceRoutes() {
	app.router.HandleFunc("/workspaces/available", app.GetAvailability).
		Methods("GET").
		Queries("floor", "{floor}").
		Queries("start", "{start:[0-9]+}").
		Queries("end", "{end:[0-9]+}")
	app.router.HandleFunc("/bulk/workspaces", app.BulkCreateWorkspaces).Methods("POST")
	app.router.HandleFunc("/workspaces", app.CreateWorkspace).Methods("POST")
	app.router.HandleFunc("/workspaces/{id}", app.GetOneWorkspace).Methods("GET")
	app.router.HandleFunc("/workspaces", app.GetAllWorkspacesByFloorId).Methods("GET").
		Queries("floor", "{floor}")
	app.router.HandleFunc("/workspaces", app.GetAllWorkspaces)
	app.router.HandleFunc("/workspaces/{id}/props", app.UpdateWorkspaceProps).Methods("PATCH")
	app.router.HandleFunc("/workspaces/{id}", app.UpdateWorkspace).Methods("PATCH")
	//app.router.HandleFunc("/workspaces/{id}", app.DeleteWorkspace).Methods("DELETE")
	app.router.HandleFunc("/assignments", app.CreateAssignments).Methods("POST")
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

func (app *App) GetAllWorkspacesByFloorId(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	floorId := queryParams["floor"][0]
	workspaces, err := app.store.WorkspaceProvider.GetAllWorkspacesByFloor(floorId)
	if err != nil {
		log.Printf("App.GetAllWorkspacesByFloorId - error getting all workspaces by floor id from provider %v", err)
		w.WriteHeader(http.StatusBadRequest)
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
		log.Printf("App.UpdateWorkspace - error updating workspace from provider %v", err)
		if strings.Contains(err.Error(), "workspace name already exists") {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedWorkspace)
}

func (app *App) UpdateWorkspaceProps(w http.ResponseWriter, r *http.Request) {
	workspaceID := mux.Vars(r)["id"]

	if workspaceID == "" {
		log.Printf("App.UpdateWorkspaceMetadata - empty workspace id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var updatedProperties model.Attrs
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("App.UpdateWorkspaceMetadata - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &updatedProperties)
	if err != nil {
		log.Printf("App.UpdateWorkspaceMetadata - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = app.store.WorkspaceProvider.UpdateWorkspaceMetadata(workspaceID, &updatedProperties)
	if err != nil {
		log.Printf("App.UpdateWorkspaceMetadata - error updating workspace from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//func (app *App) DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
//	workspaceID := mux.Vars(r)["id"]
//
//	if workspaceID == "" {
//		log.Printf("App.DeleteWorkspace - empty workspace id")
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err := app.store.WorkspaceProvider.RemoveWorkspace(workspaceID)
//	if err != nil {
//		log.Printf("App.DeleteWorkspace - error getting all workspaces from provider %v", err)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//}

func (app *App) GetAvailability(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	start := queryParams["start"][0]
	end := queryParams["end"][0]
	floorId := queryParams["floor"][0]
	if floorId == "" {
		log.Printf("App.GetAvailability - empty floor id param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	startTime, errStart := utils.TimeStampToTime(start) // Unix Timestamp
	endTime, errEnd := utils.TimeStampToTime(end)
	if errStart != nil {
		log.Printf("App.GetAvailability - empty start time param: %v", errStart)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errEnd != nil {
		log.Printf("App.GetAvailability - empty end time param: %v", errEnd)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	workspaceIds, err := app.store.WorkspaceProvider.FindAvailability(floorId, startTime, endTime)
	if err != nil {
		log.Printf("App.GetOfferingsByDateRange - error getting ids from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(workspaceIds)
}

func (app *App) CreateAssignments(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize+512)
	parseErr := r.ParseMultipartForm(MaxFileSize)
	if parseErr != nil {
		log.Println("App.CreateAssignments - failed to parse message")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		log.Println("App.CreateAssignments - expecting multipart form file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	assignmentsFile, _, err := r.FormFile("assignments")
	if err != nil {
		log.Println("App.CreateAssignments - users file is absent: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	floors, err := app.store.FloorProvider.GetAllFloors()
	if err != nil {
		log.Println("App.CreateAssignments - failed to get all floors: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	floorMap := make(map[string]string)
	for _, f := range floors {
		floorMap[f.Name] = f.ID
	}
	csvFile := csv.NewReader(assignmentsFile) // workspaceName, FloorName, UserId
	workspaces := make([]*model.Workspace, 0)
	_, _ = csvFile.Read() // skip first row
	for {
		// Read each record from csv
		record, err := csvFile.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) != 3 {
			log.Println("App.CreateAssignments - failed to parse csv file")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		workspaceName := strings.TrimSpace(record[0])
		floorName := strings.TrimSpace(record[1])
		userId := strings.TrimSpace(record[2])
		workspace := &model.Workspace{
			Name:  workspaceName,
			Floor: floorMap[floorName],
			Props: nil,
		}
		id, err := app.store.WorkspaceProvider.CreateAssignWorkspace(workspace, userId)
		if err != nil {
			log.Println("App.CreateAssignments - failed to create workspace-assignment: " + err.Error())
		}
		workspace.ID = id
		workspaces = append(workspaces, workspace)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workspaces)
}

func (app *App) BulkCreateWorkspaces(w http.ResponseWriter, r *http.Request) {
	var input model.BulkCreateWorkspacesInput
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("App.BulkCreateWorkspaces - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &input)
	if err != nil {
		log.Printf("App.BulkCreateWorkspaces - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(input)

	createdWorkspaces := make([]*model.Workspace, 0)
	for _, ws := range input.Workspaces {
		log.Println(ws)
		workspace := &model.Workspace{
			Floor:   input.FloorId,
			Name:    ws.WorkspaceName,
			Props:   ws.Props,
			Details: ws.Details,
		}
		workspaceID, err := app.store.WorkspaceProvider.UpsertWorkspace(workspace)
		if err != nil {
			log.Printf(
				"App.BulkCreateWorkspaces - failed to update details for workspace %+v with floor %s - err: %+v\n",
				ws, input.FloorId, err,
			)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(createdWorkspaces)
			return
		}
		workspace.ID = workspaceID
		createdWorkspaces = append(createdWorkspaces, workspace)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdWorkspaces)
}
