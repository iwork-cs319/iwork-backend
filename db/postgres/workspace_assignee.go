package postgres

import (
	"database/sql"
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
