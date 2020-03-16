package postgres

import (
	"database/sql"
	"log"
)

func (p PostgresDBStore) IsAssigned(id string) (bool, error) {
	sqlStatement := `SELECT workspace_id FROM workspace_assignee WHERE workspace_id=$1;`
	var returned string
	log.Printf("alive")
	row := p.database.QueryRow(sqlStatement, id)
	log.Printf("about to enter switch")
	switch err := row.Scan(&returned); err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}
