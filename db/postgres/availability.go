package postgres

import (
	"log"
	"time"
)

func (p PostgresDBStore) FindAvailability(floorId string, start time.Time, end time.Time) ([]string, error) {
	// $1 = floor-id
	// $2 = start
	// $3 = end
	log.Println(floorId, start, end)
	sqlStmtAll := `SELECT id from workspaces where floor_id=$1;`
	rows, err := p.database.Query(sqlStmtAll, floorId)
	if err != nil {
		log.Printf("PostgresDBStore.FindAvailability: %v, sqlStatement: %s\n", err, sqlStmtAll)
		return nil, err
	}
	defer rows.Close()
	allWorkspaces := make([]string, 0)
	for rows.Next() {
		var id string
		err := rows.Scan(
			&id,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.FindAvailability: %v, sqlStatement: %s\n", err, sqlStmtAll)
		}
		allWorkspaces = append(allWorkspaces, id)
	}

	sqlStmtBookings := `SELECT b.id from bookings b INNER JOIN workspaces w ON w.id=b.workspace_id
			WHERE w.floor_id=$1 AND w.locked=false AND 
			( (b.start_time <= $2 AND b.end_time >= $2) 
				OR (b.start_time >= $2 AND b.end_time <= $3) 
				OR (b.end_time >= $3 AND b.start_time <= $3) );`
	rows, err = p.database.Query(sqlStmtBookings, floorId, start, end)
	if err != nil {
		log.Printf("PostgresDBStore.FindAvailability: %v, sqlStatement: %s\n", err, sqlStmtBookings)
		return nil, err
	}
	defer rows.Close()
	bookedWorkspaces := make([]string, 0)
	for rows.Next() {
		var id string
		err := rows.Scan(
			&id,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.FindAvailability: %v, sqlStatement: %s\n", err, sqlStmtBookings)
		}
		bookedWorkspaces = append(bookedWorkspaces, id)
	}

	availableWorkspaces := difference(allWorkspaces, bookedWorkspaces)
	return availableWorkspaces, nil
}

// Set Difference: A - B
func difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}
