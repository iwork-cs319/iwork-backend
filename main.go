package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type BasicMessage struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type workspace struct {
	ID          string `json:"ID"`
	Name       string `json:"Name"`
	Availability bool `json:"Availability"`
}

type allWorkspaces []workspace

var workspaces = allWorkspaces {
	{
		ID:          "1",
		Name:       "Workspace #001",
		Availability: true,
	},
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(BasicMessage{
		"1",
		"Hello World!",
	})
}

func createWorkspace(w http.ResponseWriter, r *http.Request) {
	var newWorkspace workspace
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the workspace title and availability only in order to update")
	}

	json.Unmarshal(reqBody, &newWorkspace)
	workspaces = append(workspaces, newWorkspace)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newWorkspace)
}

func getOneWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceID := mux.Vars(r)["id"]

	for _, singleWorkspace := range workspaces {
		if singleWorkspace.ID == workspaceID {
			json.NewEncoder(w).Encode(singleWorkspace)
		}
	}
}

func getAllWorkspaces(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(workspaces)
}

func updateWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceID := mux.Vars(r)["id"]
	var updatedWorkspace workspace

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the workspace title and availability only in order to update")
	}
	json.Unmarshal(reqBody, &updatedWorkspace)

	for i, singleWorkspace := range workspaces {
		if singleWorkspace.ID == workspaceID {
			singleWorkspace.Name = updatedWorkspace.Name
			singleWorkspace.Availability = updatedWorkspace.Availability
			workspaces = append(workspaces[:i], singleWorkspace)
			json.NewEncoder(w).Encode(singleWorkspace)
		}
	}
}

func deleteWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceID := mux.Vars(r)["id"]

	for i, singleWorkspace := range workspaces {
		if singleWorkspace.ID == workspaceID {
			workspaces = append(workspaces[:i], workspaces[i+1:]...)
			fmt.Fprintf(w, "The workspace with ID %v has been deleted successfully", workspaceID)
		}
	}
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You've reached the home endpoint!")
}
