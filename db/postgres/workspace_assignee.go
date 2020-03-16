package postgres

import (
	"database/sql"
)

func (p PostgresDBStore) IsAssigned(id string) (bool, error) {
	sqlStatement := `SELECT workspace_id FROM workspace_assignee WHERE workspace_id=$1;`
	var returned string;
	row := p.database.QueryRow(sqlStatement, id)
	switch err := row.Scan(&returned); err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}
