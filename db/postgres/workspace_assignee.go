package postgres

import (
	"database/sql"
	"go-api/model"
	"go-api/utils"
	"log"
	"time"
)

func (p PostgresDBStore) IsAssigned(id string, startTime time.Time, endTime time.Time) (bool, error) {
	// assigned start <= booking start & assigned end either null or >= booking end
	sqlStatement := `SELECT workspace_id FROM workspace_assignee
					 	WHERE workspace_id=$1 AND (
					 			(start_time <= $2 AND end_time >=$2) OR
					 			(start_time <= $3 AND end_time >=$3) OR
					 			(start_time >= $2 AND end_time <=$3) OR 
					 			(start_time <= $2 AND end_time >=$3)
								)`
	var returned string
	row := p.database.QueryRow(sqlStatement, id, startTime, endTime)
	switch err := row.Scan(&returned); err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

func (p PostgresDBStore) IsFullyAssigned(id string, startTime time.Time, endTime time.Time) (bool, error) {
	// assigned start <= booking start & assigned end either null or >= booking end
	sqlStatement := `SELECT workspace_id FROM workspace_assignee
					 	WHERE workspace_id=$1 AND (
					 	    (start_time <= $1 AND end_time >=$3));`
	var returned string
	row := p.database.QueryRow(sqlStatement, id, startTime, endTime)
	switch err := row.Scan(&returned); err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

func (p PostgresDBStore) GetExpiredAssignments(since time.Time) ([]*model.Assignment, error) {
	rows, err := p.database.Query(
		`SELECT id, workspace_id, user_id, start_time, end_time FROM workspace_assignee WHERE end_time < $1`,
		since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	assignments := make([]*model.Assignment, 0)
	for rows.Next() {
		var assignment model.Assignment
		err := rows.Scan(
			&assignment.ID,
			&assignment.WorkspaceID,
			&assignment.UserID,
			&assignment.StartDate,
			&assignment.EndDate,
		)
		if err != nil {
			log.Printf("PostgresDBStore.GetExpiredAssignments: %v\n", err)
		}
		if assignment.UserID != utils.EmptyUserUUID {
			assignments = append(assignments, &assignment)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return assignments, nil
}
