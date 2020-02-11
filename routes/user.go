package routes

import (
"encoding/json"
"github.com/gorilla/mux"
//"go-api/model"
//"io/ioutil"
"log"
"net/http"
)

func (app *App) RegisterUserRoutes() {
	//app.router.HandleFunc("/users", app.CreateUser).Methods("POST")
	app.router.HandleFunc("/users/{id}", app.GetOneUser).Methods("GET")
	//app.router.HandleFunc("/users/workspaces/{workspace_id}", app.GetUsersByWorkspaceID).Methods("GET")
	app.router.HandleFunc("/users", app.GetAllUsers).Methods("GET")
	//app.router.HandleFunc("/users/{id}", app.UpdateUser).Methods("PATCH")
	//app.router.HandleFunc("/users/{id}", app.RemoveUser).Methods("DELETE")
}

//func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
//	var newUser model.User
//	reqBody, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Printf("App.CreateUser - error reading request body %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err = json.Unmarshal(reqBody, &newUser)
//	if err != nil {
//		log.Printf("App.CreateUser - error unmarshaling request body %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err = app.store.UserProvider.CreateUser(&newUser)
//	if err != nil {
//		log.Printf("App.CreateUser - error creating user %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusCreated)
//	json.NewEncoder(w).Encode(newUser)
//}

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

