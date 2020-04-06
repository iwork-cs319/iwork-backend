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
	"time"
)

func (app *App) RegisterWorkspaceRoutes() {
	app.router.HandleFunc("/workspaces/available", app.GetAvailability).
		Methods("GET").
		Queries("floor", "{floor}").
		Queries("start", "{start:[0-9]+}").
		Queries("end", "{end:[0-9]+}")
	app.router.HandleFunc("/workspaces/available", app.GetAllFloorsAvailability).
		Methods("GET").
		Queries("start", "{start:[0-9]+}").
		Queries("end", "{end:[0-9]+}")
	app.router.HandleFunc("/workspaces/bulk/available", app.GetBulkAvailability).
		Methods("GET").
		Queries("start", "{start:[0-9]+}").
		Queries("end", "{end:[0-9]+}")
	app.router.HandleFunc("/workspaces/bulk/countavailable", app.GetBulkCountAvailability).
		Methods("GET").
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
	//app.router.HandleFunc("/workspaces/store/available", app.GetAvailabilityYesterday).Methods("GET")
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
		log.Printf("App.FindAvailability - error getting ids from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(workspaceIds)
}

func (app *App) GetAllFloorsAvailability(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	start := queryParams["start"][0]
	end := queryParams["end"][0]
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
	floorIDs, err := app.store.FloorProvider.GetAllFloorIDs()
	if err != nil {
		log.Printf("App.GetAllFloorsAvailability - error getting floor_id's from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	allWorkspaceIDs := make([]string, 0)
	for _, f := range floorIDs {
		workspaceIDs, err := app.store.WorkspaceProvider.FindAvailability(f, startTime, endTime)
		if err != nil {
			log.Printf("App.GetAllFloorsAvailability - error getting ids from provider %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		allWorkspaceIDs = append(allWorkspaceIDs, workspaceIDs...)
	}
	json.NewEncoder(w).Encode(allWorkspaceIDs)
}

//type DateWorkspaceStat struct {
//	DayWorkspaceStat  map[string]WorkspaceStat `json:"workspace_stat"`
//	CountDayAvailable int                      `json:"countd_available"`
//	CountDayTotal     int                      `json:"countd_total"`
//}

func (app *App) GetBulkCountAvailability(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	start := queryParams["start"][0]
	end := queryParams["end"][0]
	startTime, errStart := utils.TimeStampToTime(start) // Unix Timestamp
	endTime, errEnd := utils.TimeStampToTime(end)
	if errStart != nil {
		log.Printf("App.GetBulkAvailability - empty start time param: %v", errStart)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errEnd != nil {
		log.Printf("App.GetBulkAvailability - empty end time param: %v", errEnd)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if endTime.Before(startTime) {
		log.Printf("App.GetBulkAvailability - end time before start time: %v", errStart)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loc, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		log.Printf("App.GetBulkAvailability - location failed with: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Truncate start time to beginning of the day
	startT := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, loc)
	endT := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 23, 59, 59, 0, loc)
	finalEnd := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, loc)
	// Ensure end time is end of it's day
	// For loop for each day, collecting the information into a dict
	allDaysDict := make(map[string]map[string]WorkspaceCount)
	for s, e := startT, endT; s.Before(finalEnd); s, e = s.AddDate(0, 0, 1), e.AddDate(0, 0, 1) {
		workspacesDict, err := app.getBulkCountAvailabilities(s, e)
		if err != nil {
			log.Printf("App.GetBulkCountAvailability - error getting BulkAvailabilities %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Create dict id with d
		allDaysDict[s.Format("02.01.2006")] = workspacesDict
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allDaysDict)
	return
}

func (app *App) GetBulkAvailability(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	start := queryParams["start"][0]
	end := queryParams["end"][0]
	startTime, errStart := utils.TimeStampToTime(start) // Unix Timestamp
	endTime, errEnd := utils.TimeStampToTime(end)
	if errStart != nil {
		log.Printf("App.GetBulkAvailability - empty start time param: %v", errStart)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errEnd != nil {
		log.Printf("App.GetBulkAvailability - empty end time param: %v", errEnd)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if endTime.Before(startTime) {
		log.Printf("App.GetBulkAvailability - end time before start time: %v", errStart)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Truncate start time to beginning of the day
	startT := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())  // TODO: Check location
	endT := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 23, 59, 59, 0, startTime.Location()) // TODO: Check location
	finalEnd := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, startTime.Location())   // TODO: Check location
	// Ensure end time is end of it's day
	// For loop for each day, collecting the information into a dict
	allDaysDict := make(map[string]map[string]WorkspaceStat)
	for s, e := startT, endT; s.Before(finalEnd); s, e = s.AddDate(0, 0, 1), e.AddDate(0, 0, 1) {
		workspacesDict, err := app.getBulkAvailabilities(s, e)
		if err != nil {
			log.Printf("App.GetBulkAvailability - error getting BulkAvailabilities %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Create dict id with d
		allDaysDict[s.Format("02.01.2006")] = workspacesDict
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allDaysDict)
	return
}

type WorkspaceStat struct {
	WorkspaceCount
	WorkspaceIDs []string `json:"workspace_ids"`
}

type WorkspaceCount struct {
	CountAvailable int `json:"count_available"`
	CountFloor     int `json:"count_floor"`
}

//func (app *App) GetAvailabilityYesterday(w http.ResponseWriter, r *http.Request) {
//	// Gets called the day after to store the availabilities of the previous day.
//	yesterday, yesterdayEnd, err := utils.TimeYesterday()
//	if err != nil {
//		log.Printf("App.GetAvailabilityYesterday - error getting time %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	workspacesDict, err := app.getBulkAvailabilities(yesterday, yesterdayEnd)
//	if workspacesDict == nil {
//		log.Printf("App.GetAvailabilityYesterday - error getting BulkAvailabilities %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(workspacesDict)
//	return
//}

func (app *App) getBulkCountAvailabilities(start time.Time, end time.Time) (map[string]WorkspaceCount, error) {
	// Find all availabilities for the previous day
	floorIDs, err := app.store.FloorProvider.GetAllFloorIDs()
	if err != nil {
		log.Printf("App.getBulkAvailabilities - error getting floor_id's from provider %v", err)
		return nil, err
	}
	workspacesDict := make(map[string]WorkspaceCount)
	for _, f := range floorIDs {
		workspaceIDs, err := app.store.WorkspaceProvider.FindAvailability(f, start, end)
		if err != nil {
			log.Printf("App.getBulkAvailabilities - error getting ids from provider %v", err)
			return nil, err
		}
		count := len(workspaceIDs)
		countWorkspacesOnFloor, err := app.store.WorkspaceProvider.CountWorkspacesByFloor(f)
		// Create a map[string]WorkspaceStat for floorid -> ...
		floorStat := WorkspaceCount{
			CountAvailable: count,
			CountFloor:     countWorkspacesOnFloor,
		}
		workspacesDict[f] = floorStat
		//allWorkspaceIDs = append(allWorkspaceIDs, workspaceIDs...)
	}
	return workspacesDict, nil
}

func (app *App) getBulkAvailabilities(start time.Time, end time.Time) (map[string]WorkspaceStat, error) {
	// Find all availabilities for the previous day
	floorIDs, err := app.store.FloorProvider.GetAllFloorIDs()
	if err != nil {
		log.Printf("App.getBulkAvailabilities - error getting floor_id's from provider %v", err)
		return nil, err
	}
	workspacesDict := make(map[string]WorkspaceStat)
	for _, f := range floorIDs {
		workspaceIDs, err := app.store.WorkspaceProvider.FindAvailability(f, start, end)
		if err != nil {
			log.Printf("App.getBulkAvailabilities - error getting ids from provider %v", err)
			return nil, err
		}
		count := len(workspaceIDs)
		countWorkspacesOnFloor, err := app.store.WorkspaceProvider.CountWorkspacesByFloor(f)
		// Create a map[string]WorkspaceStat for floorid -> ...
		wC := WorkspaceCount{
			CountAvailable: count,
			CountFloor:     countWorkspacesOnFloor,
		}
		floorStat := WorkspaceStat{
			WorkspaceCount: wC,
			WorkspaceIDs:   workspaceIDs,
		}
		workspacesDict[f] = floorStat
		//allWorkspaceIDs = append(allWorkspaceIDs, workspaceIDs...)
	}
	return workspacesDict, nil
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
			log.Printf("App.CreateAssignments - failed to create workspace-assignment for w:%s u:%s err:%s ", workspace, userId, err.Error())
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
	errors := make([]*model.BulkCreateWorkspaceError, 0)
	createdWorkspaces := make([]*model.CreateWorkspaceInput, 0)
	for _, ws := range input.Workspaces {
		workspace := &model.Workspace{
			ID:      ws.WorkspaceId,
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
			errors = append(
				errors,
				&model.BulkCreateWorkspaceError{WorkspaceName: ws.WorkspaceName, Message: err.Error()},
			)
			continue
		}
		if ws.WorkspaceId == "" && ws.UserId == "" {
			// New workspace with no user
			// Create a default offering
			_, err = app.store.OfferingProvider.CreateDefaultOffering(&model.Offering{
				UserID:      utils.EmptyUserUUID,
				WorkspaceID: workspaceID,
				StartDate:   time.Now(),
				CreatedBy:   utils.EmptyUserUUID,
			})
			if err != nil {
				log.Printf(
					"App.BulkCreateWorkspaces - failed to create a default offering for workspace %+v with floor %s - err: %+v\n",
					ws, input.FloorId, err,
				)
				errors = append(
					errors,
					&model.BulkCreateWorkspaceError{WorkspaceName: ws.WorkspaceName, Message: err.Error()},
				)
				continue
			}
		} else if ws.UserId != "" {
			// Create an assignment
			err = app.store.WorkspaceProvider.CreateAssignment(ws.UserId, workspaceID)
			if err != nil {
				log.Printf(
					"App.BulkCreateWorkspaces - failed to create an assignment for workspace %+v with floor %s - err: %+v\n",
					ws, input.FloorId, err,
				)
				errors = append(
					errors,
					&model.BulkCreateWorkspaceError{WorkspaceName: ws.WorkspaceName, Message: err.Error()},
				)
				continue
			}
		}
		ws.WorkspaceId = workspaceID
		createdWorkspaces = append(createdWorkspaces, ws)
	}
	if len(createdWorkspaces) > 0 {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	ret := struct {
		Data   []*model.CreateWorkspaceInput     `json:"data"`
		Errors []*model.BulkCreateWorkspaceError `json:"errors"`
	}{
		createdWorkspaces,
		errors,
	}
	json.NewEncoder(w).Encode(ret)
}
