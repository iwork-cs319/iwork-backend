package postgres

import (
	"go-api/model"
	"go-api/utils"
	"log"
	"time"
)

func (p PostgresDBStore) GetOneUser(id string) (*model.User, error) {
	sqlStatement := `SELECT id, name, email, department, is_admin FROM users WHERE id=$1;`
	var user model.User
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Department,
		&user.IsAdmin,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p PostgresDBStore) GetAllUsers() ([]*model.User, error) {
	sqlStatement := `SELECT id, name, email, department, is_admin FROM users WHERE deleted=FALSE;`
	return p.queryMultipleUsers(sqlStatement)
}

func (p PostgresDBStore) CreateUser(user *model.User) error {
	sqlStatement :=
		`INSERT INTO users(id, name, email, department, is_admin) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement,
		user.ID,
		user.Name,
		user.Email,
		user.Department,
		user.IsAdmin,
	).Scan(&id)
	if err != nil {
		return err
	}
	if id != user.ID {
		return CreateError
	}
	return nil
}

func (p PostgresDBStore) GetAssignedUsers(start, end time.Time) ([]*model.UserAssignment, error) {
	getOfferingsStmt := `SELECT id, user_id, workspace_id, start_time, end_time, cancelled, created_by from offerings 
							WHERE cancelled=FALSE AND ((start_time <= $1 AND end_time >= $2) OR 
									(start_time >= $1 AND end_time <= $2) OR 
									(start_time <= $1 AND end_time >= $1) OR
									(start_time <= $2 AND end_time >= $2))`
	offerings, err := p.queryMultipleOfferings(getOfferingsStmt, start, end)
	if err != nil {
		return nil, err
	}

	sqlStatement := `SELECT users.id, name, email, department, is_admin, wa.workspace_id FROM users 
							INNER JOIN workspace_assignee wa ON users.id = wa.user_id
							WHERE users.deleted=FALSE AND wa.start_time <= $1 AND (wa.end_time >= $2 OR end_time IS NULL)`
	rows, err := p.database.Query(sqlStatement, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	assignedUsers := make([]*model.UserAssignment, 0)
	for rows.Next() {
		var user model.UserAssignment
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Department,
			&user.IsAdmin,
			&user.WorkspaceId,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.GetAssignedUsers: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		assignedUsers = append(assignedUsers, &user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	users := make([]*model.UserAssignment, 0)
	for _, u := range assignedUsers {
		foundOfferingByUser := false
		for _, o := range offerings {
			if o.UserID == u.ID {
				foundOfferingByUser = true
				break
			}
		}
		if !foundOfferingByUser {
			users = append(users, u)
		}
	}

	return users, nil
}

func (p PostgresDBStore) GetAssignedUsersByTime(timestamp time.Time) ([]*model.UserAssignment, error) {
	sqlStatement := `SELECT users.id, name, email, department, is_admin, wa.workspace_id FROM users 
							INNER JOIN workspace_assignee wa ON users.id = wa.user_id
							WHERE wa.start_time <= $1 AND (wa.end_time >= $1 OR end_time IS NULL)`
	rows, err := p.database.Query(sqlStatement, timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	assignedUsers := make([]*model.UserAssignment, 0)
	for rows.Next() {
		var user model.UserAssignment
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Department,
			&user.IsAdmin,
			&user.WorkspaceId,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.GetAssignedUsersByTime: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		assignedUsers = append(assignedUsers, &user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return assignedUsers, nil
}

//func (p PostgresDBStore) UpdateUser(id string, user *model.User) error {
//	sqlStatement :=
//		`UPDATE users
//				SET name = $2, department = $3, is_admin = $4
//				WHERE id = $1
//				RETURNING id;`
//	var _id string
//	err := p.database.QueryRow(sqlStatement,
//		id,
//		user.Name,
//		user.Department,
//		user.IsAdmin,
//	).Scan(&_id)
//	if err != nil {
//		return err
//	}
//	if _id != id {
//		return CreateError
//	}
//	return nil
//}
//
//func (p PostgresDBStore) RemoveUser(id string) error {
//	sqlStatement :=
//		`DELETE FROM users
//				WHERE id = $1
//				RETURNING id;`
//	var _id string
//	err := p.database.QueryRow(sqlStatement,
//		id,
//	).Scan(&_id)
//	if err != nil {
//		return err
//	}
//	if _id != id {
//		return CreateError
//	}
//	return nil
//}

func (p PostgresDBStore) queryMultipleUsers(sqlStatement string, args ...interface{}) ([]*model.User, error) {
	rows, err := p.database.Query(sqlStatement, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*model.User, 0)
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Department,
			&user.IsAdmin,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.queryMultipleOfferings: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		if user.ID != utils.EmptyUserUUID {
			users = append(users, &user)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}
