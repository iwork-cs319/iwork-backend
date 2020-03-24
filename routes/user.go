package routes

import (
	"encoding/csv"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/model"
	"go-api/utils"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (app *App) RegisterUserRoutes() {
	app.router.HandleFunc("/users/assigned", app.GetAllAssignedUsers).Methods("GET").Queries("start", "{start:[0-9]+}").Queries("end", "{end:[0-9]+}")
	app.router.HandleFunc("/users", app.CreateUsers).Methods("POST")
	app.router.HandleFunc("/users/{id}", app.GetOneUser).Methods("GET")
	//app.router.HandleFunc("/users/workspaces/{workspace_id}", app.GetUsersByWorkspaceID).Methods("GET")
	app.router.HandleFunc("/users", app.GetAllUsers).Methods("GET")
	//app.router.HandleFunc("/users/{id}", app.UpdateUser).Methods("PATCH")
	//app.router.HandleFunc("/users/{id}", app.RemoveUser).Methods("DELETE")
}

func (app *App) CreateUsers(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize+512)
	parseErr := r.ParseMultipartForm(MaxFileSize)
	if parseErr != nil {
		log.Println("App.CreateUsers - failed to parse message")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		log.Println("App.CreateUsers - expecting multipart form file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usersFile, _, err := r.FormFile("users")
	if err != nil {
		log.Println("App.CreateUsers - users file is absent: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	csvFile := csv.NewReader(usersFile) // email, name, id, department, isAdmin
	users := make([]*model.User, 0)
	_, _ = csvFile.Read() // skip heading row
	for {
		// Read each record from csv
		record, err := csvFile.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) != 5 {
			log.Println("App.CreateUsers - failed to parse csv file")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		b, err := strconv.ParseBool(record[4])
		if err != nil {
			b = false
		}
		email := strings.TrimSpace(record[0])
		name := strings.TrimSpace(record[1])
		userId := strings.TrimSpace(record[2])
		department := strings.TrimSpace(record[3])
		user := &model.User{
			Email:      email,
			Name:       name,
			ID:         userId,
			Department: department,
			IsAdmin:    b,
		}
		err = app.store.UserProvider.CreateUser(user)
		if err != nil {
			log.Printf("App.CreateUsers - error creating user %+v\n", err)
		} else {
			users = append(users, user)
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(users)
}

func (app *App) GetOneUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	if userID == "" {
		log.Printf("App.GetOneUser - empty user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := app.store.UserProvider.GetOneUser(userID)
	if err != nil {
		log.Printf("App.GetOneUser - error getting user from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (app *App) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.store.UserProvider.GetAllUsers()
	if err != nil {
		log.Printf("App.GetAllUsers - error getting all users from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (app *App) GetAllAssignedUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	start := queryParams["start"][0]
	end := queryParams["end"][0]
	startTime, errStart := utils.TimeStampToTime(start) // Unix Timestamp
	endTime, errEnd := utils.TimeStampToTime(end)
	if errStart != nil {
		log.Printf("App.GetAllAssignedUsers - empty start time param: %v", errStart)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errEnd != nil {
		log.Printf("App.GetAllAssignedUsers - empty end time param: %v", errEnd)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	users, err := app.store.UserProvider.GetAssignedUsers(startTime, endTime)
	if err != nil {
		log.Printf("App.GetAllAssignedUsers - error getting all users from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

//func (app *App) GetUsersByWorkspaceID(w http.ResponseWriter, r *http.Request) {
//	workspaceID := mux.Vars(r)["workspace_id"]
//
//	if workspaceID == "" {
//		log.Printf("App.GetOneUser - empty user id")
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	users, err := app.store.UserProvider.GetUsersByWorkspaceID(workspaceID)
//	if err != nil {
//		log.Printf("App.GetUsersByWorkspaceID - error getting users by workspaceID from provider %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	json.NewEncoder(w).Encode(users)
//}

//func (app *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
//	userID := mux.Vars(r)["id"]
//
//	if userID == "" {
//		log.Printf("App.UpdateUser - empty user id")
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	var updatedUser model.User
//	reqBody, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Printf("App.UpdateUser - error reading request body %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	err = json.Unmarshal(reqBody, &updatedUser)
//	if err != nil {
//		log.Printf("App.UpdateUser - error unmarshaling request body %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err = app.store.UserProvider.UpdateUser(userID, &updatedUser)
//	if err != nil {
//		log.Printf("App.UpdateUser - error getting all users from provider %v", err)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//}

//func (app *App) RemoveUser(w http.ResponseWriter, r *http.Request) {
//	userID := mux.Vars(r)["id"]
//
//	if userID == "" {
//		log.Printf("App.RemoveUser - empty user id")
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err := app.store.UserProvider.RemoveUser(userID)
//	if err != nil {
//		log.Printf("App.RemoveUser - error getting all users from provider %v", err)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//}
