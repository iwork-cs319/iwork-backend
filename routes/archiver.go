package routes

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func (app *App) RegisterArchiverRoutes() {
	app.router.HandleFunc("/archiver", app.Archive).Methods("POST")
}

func (app *App) Archive(w http.ResponseWriter, r *http.Request) {
	writer := new(bytes.Buffer)
	archiveDate := time.Now().UTC().Format(time.RFC822Z)
	archiveFileName := fmt.Sprintf("Test--%s", archiveDate)
	writer.WriteString(archiveFileName) // Write header
	writer.WriteString("=Bookings(id,workspace_id,user_id,created_by,start_date,end_date,cancelled")
	bookings, err := app.store.BookingProvider.GetExpiredBookings()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, b := range bookings {
		writer.WriteString(fmt.Sprintf(
			"%s,%s,%s,%s,%s,%s,%s",
			b.ID, b.WorkspaceID, b.UserID, b.CreatedBy,
			b.StartDate.Format(time.RFC3339), b.EndDate.Format(time.RFC3339), b.Cancelled,
		))
	}
	err = app.gDrive.UploadArchiveDataFile(archiveFileName, writer)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
