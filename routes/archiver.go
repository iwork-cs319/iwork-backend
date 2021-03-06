package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (app *App) RegisterArchiverRoutes() {
	app.router.HandleFunc("/archiver", app.Archive).Methods("POST")
}

func (app *App) Archive(w http.ResponseWriter, r *http.Request) {
	writer := new(bytes.Buffer)

	now := time.Now()
	archiveDateTimestamp := fmt.Sprintf("%d", now.Unix())
	archiveFileName := fmt.Sprintf("Test_%s.archive", archiveDateTimestamp)
	writer.WriteString(fmt.Sprintf("%s\n", archiveFileName)) // Write file name

	if err := app.WriteOfferings(writer, now); err != nil {
		WriteLine(writer, "-- error writing offerings")
		log.Println("App.Archive.Offerings:", err)
	}

	if err := app.WriteBookings(writer, now); err != nil {
		WriteLine(writer, "-- error writing bookings")
		log.Println("App.Archive.Bookings:", err)
	}

	if err := app.WriteWorkspaceAssignments(writer, now); err != nil {
		WriteLine(writer, "-- error writing assignments")
		log.Println("App.Archive.Assignments:", err)
	}

	if err := app.WriteWorkspaces(writer); err != nil {
		WriteLine(writer, "-- error writing workspaces")
		log.Println("App.Archive.Workspaces:", err)
	}

	if err := app.WriteFloors(writer); err != nil {
		WriteLine(writer, "-- error writing floors")
		log.Println("App.Archive.Floors:", err)
	}

	if err := app.gDrive.UploadArchiveDataFile(archiveFileName, writer); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func WriteLine(w *bytes.Buffer, header string) {
	w.WriteString(fmt.Sprintf("%s\n", header))
}

func (app *App) WriteBookings(w *bytes.Buffer, now time.Time) error {
	bookings, err := app.store.BookingProvider.GetExpiredBookings(now)
	if err != nil {
		return err
	}
	WriteLine(w, "== Bookings(id|~|workspace_id|~|user_id|~|created_by|~|start_date|~|end_date|~|cancelled)")
	ids := make([]string, 0)
	for _, b := range bookings {
		WriteLine(w, fmt.Sprintf(
			"%s|~|%s|~|%s|~|%s|~|%s|~|%s|~|%t",
			b.ID, b.WorkspaceID, b.UserID, b.CreatedBy,
			b.StartDate.Format(time.RFC3339), b.EndDate.Format(time.RFC3339), b.Cancelled,
		))
		ids = append(ids, b.ID)
	}
	return app.store.BookingProvider.DeleteBookings(ids)
}

func (app *App) WriteOfferings(w *bytes.Buffer, now time.Time) error {
	offerings, err := app.store.OfferingProvider.GetExpiredOfferings(now)
	if err != nil {
		return err
	}
	WriteLine(w, "== Offerings(id|~|workspace_id|~|user_id|~|created_by|~|start_date|~|end_date|~|cancelled)")
	ids := make([]string, 0)
	for _, o := range offerings {
		WriteLine(w, fmt.Sprintf(
			"%s|~|%s|~|%s|~|%s|~|%s|~|%s|~|%t",
			o.ID, o.WorkspaceID, o.UserID, o.CreatedBy,
			o.StartDate.Format(time.RFC3339), o.EndDate.Format(time.RFC3339), o.Cancelled,
		))
		ids = append(ids, o.ID)
	}
	return app.store.OfferingProvider.DeleteOfferings(ids)
}

func (app *App) WriteFloors(w *bytes.Buffer) error {
	floors, err := app.store.FloorProvider.GetDeletedFloors()
	if err != nil {
		return err
	}
	WriteLine(w, "== Floors(id|~|name|~|address|~|download_url)")
	ids := make([]string, 0)
	for _, f := range floors {
		WriteLine(w, fmt.Sprintf(
			"%s|~|%s|~|%s|~|%s",
			f.ID, f.Name, f.Address, f.DownloadURL,
		))
		ids = append(ids, f.ID)
	}
	return app.store.FloorProvider.DeleteFloors(ids)
}

func (app *App) WriteWorkspaces(w *bytes.Buffer) error {
	workspaces, err := app.store.WorkspaceProvider.GetDeletedWorkspaces()
	if err != nil {
		return err
	}
	WriteLine(w, "== Workspaces(id|~|name|~|floor_id|~|details|~|props)")
	ids := make([]string, 0)
	for _, workspace := range workspaces {
		props, _ := json.Marshal(workspace.Props)
		WriteLine(w, fmt.Sprintf(
			"%s|~|%s|~|%s|~|%s|~|%s",
			workspace.ID, workspace.Name, workspace.Floor, workspace.Details, string(props),
		))
		ids = append(ids, workspace.ID)
	}
	return app.store.WorkspaceProvider.DeleteWorkspaces(ids)
}

func (app *App) WriteWorkspaceAssignments(w *bytes.Buffer, now time.Time) error {
	assignments, err := app.store.AssigneeProvider.GetExpiredAssignments(now)
	if err != nil {
		return err
	}
	WriteLine(w, "== Assignments(id|~|workspace_id|~|user_id|~|start_time|~|end_time)")
	ids := make([]string, 0)
	for _, a := range assignments {
		WriteLine(w, fmt.Sprintf(
			"%s|~|%s|~|%s|~|%s|~|%s",
			a.ID, a.WorkspaceID, a.UserID,
			a.StartDate.Format(time.RFC3339), a.EndDate.Format(time.RFC3339),
		))
		ids = append(ids, a.ID)
	}
	return app.store.AssigneeProvider.DeleteAssignments(ids)
}
