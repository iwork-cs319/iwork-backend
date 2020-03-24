package postgres

import (
	"log"
	"time"
)

func (p PostgresDBStore) FindAvailability(floorId string, start time.Time, end time.Time) ([]string, error) {
	// $1 = floor-id
	// $2 = start
	// $3 = end

	sqlStmtBookings := `SELECT b.workspace_id from bookings b INNER JOIN workspaces w ON w.id=b.workspace_id
			WHERE w.floor_id=$1 AND b.cancelled=false AND
			( (b.start_time <= $2 AND b.end_time >= $2) 
				OR (b.start_time >= $2 AND b.end_time <= $3) 
				OR (b.end_time >= $3 AND b.start_time <= $3) );`
	rows, err := p.database.Query(sqlStmtBookings, floorId, start, end)
	if err != nil {
		log.Printf("PostgresDBStore.FindAvailability.getBookings: %v, sqlStatement: %s\n", err, sqlStmtBookings)
		return nil, err
	}
	defer rows.Close()
	bookedWorkspaces := make(map[string]bool)
	for rows.Next() {
		var id string
		err := rows.Scan(
			&id,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.FindAvailability.getBookings: %v, sqlStatement: %s\n", err, sqlStmtBookings)
		}
		bookedWorkspaces[id] = true
	}

	sqlStmtOfferings := `SELECT o.workspace_id from offerings o INNER JOIN workspaces w ON w.id=o.workspace_id
			WHERE w.floor_id=$1 AND o.cancelled=false AND
			( (o.start_time <= $2 AND o.end_time >= $3) );`
	rows, err = p.database.Query(sqlStmtOfferings, floorId, start, end)
	if err != nil {
		log.Printf("PostgresDBStore.FindAvailability.getOfferings: %v, sqlStatement: %s\n", err, sqlStmtOfferings)
		return nil, err
	}
	defer rows.Close()
	offeredWorkspaces := make(map[string]bool)
	for rows.Next() {
		var id string
		err := rows.Scan(
			&id,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.FindAvailability.getOfferings: %v, sqlStatement: %s\n", err, sqlStmtOfferings)
		}
		offeredWorkspaces[id] = true
	}

	sqlStmtAssigned := `SELECT wa.workspace_id from workspace_assignee wa INNER JOIN workspaces w ON w.id=wa.workspace_id
			WHERE w.floor_id=$1 AND
			( (wa.start_time <= $2 AND wa.end_time >= $2) 
				OR (wa.start_time >= $2 AND wa.end_time <= $3) 
				OR (wa.end_time >= $3 AND wa.start_time <= $3)
			    OR (wa.end_time IS NULL AND wa.start_time < $2)
			    OR (wa.end_time IS NULL AND wa.start_time > $2 AND wa.start_time < $3));`
	rows, err = p.database.Query(sqlStmtAssigned, floorId, start, end)
	if err != nil {
		log.Printf("PostgresDBStore.FindAvailability.getAssignee: %v, sqlStatement: %s\n", err, sqlStmtAssigned)
		return nil, err
	}
	defer rows.Close()
	assignedWorkspaces := make(map[string]bool)
	for rows.Next() {
		var id string
		err := rows.Scan(
			&id,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.FindAvailability.getAssignee: %v, sqlStatement: %s\n", err, sqlStmtAssigned)
		}
		assignedWorkspaces[id] = true
	}

	sqlStmtAll := `SELECT id from workspaces where floor_id=$1;`
	rows, err = p.database.Query(sqlStmtAll, floorId)
	if err != nil {
		log.Printf("PostgresDBStore.FindAvailability: %v, sqlStatement: %s\n", err, sqlStmtAll)
		return nil, err
	}
	defer rows.Close()
	availableWorkspaces := make([]string, 0)
	for rows.Next() {
		var id string
		err := rows.Scan(
			&id,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.FindAvailability: %v, sqlStatement: %s\n", err, sqlStmtAll)
		}
		_, offered := offeredWorkspaces[id]
		_, booked := bookedWorkspaces[id]
		_, assigned := assignedWorkspaces[id]
		if assigned && offered && !booked {
			availableWorkspaces = append(availableWorkspaces, id)
		} else if !assigned && !booked {
			availableWorkspaces = append(availableWorkspaces, id)
		}
	}

	return availableWorkspaces, nil
}
