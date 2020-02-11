package postgres

import (
	"go-api/model"
	"log"
)

func (p PostgresDBStore) GetOneUser(id string) (*model.User, error) {
	sqlStatement := `SELECT id, name, department, is_admin FROM users WHERE id=$1;`
	var user model.User
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Department,
		&user.IsAdmin,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p PostgresDBStore) GetAllUsers() ([]*model.User, error) {
	sqlStatement := `SELECT id, name, department, is_admin FROM users;`
	return p.queryMultipleUsers(sqlStatement)
}

//func (p PostgresDBStore) CreateUser(user *model.User) error {
//	sqlStatement :=
//		`INSERT INTO users(id, name, department, is_admin) VALUES ($1, $2, $3, $4) RETURNING id`
//	var id string
//	err := p.database.QueryRow(sqlStatement,
//		user.ID,
//		user.Name,
//		user.Department,
//		user.IsAdmin,
//	).Scan(&id)
//	if err != nil {
//		return err
//	}
//	if id != user.ID {
//		return CreateError
//	}
//	return nil
//}

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
			&user.Department,
			&user.IsAdmin,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.queryMultipleOfferings: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		users = append(users, &user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}
