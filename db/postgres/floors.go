package postgres

import (
	"go-api/model"
	"log"
)

func (p PostgresDBStore) GetOneFloor(id string) (*model.Floor, error) {
	sqlStatement := `SELECT id, name FROM floors WHERE id=$1;`
	var floor model.Floor
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&floor.ID,
		&floor.Name,
	)
	if err != nil {
		return nil, err
	}
	return &floor, nil
}

func (p PostgresDBStore) GetAllFloors() ([]*model.Floor, error) {
	sqlStatement := `SELECT id, name FROM floors;`
	return p.queryMultipleFloors(sqlStatement)
}

//func (p PostgresDBStore) CreateFloor(floor *model.Floor) error {
//	sqlStatement :=
//		`INSERT INTO floors(id, name) VALUES ($1, $2) RETURNING id`
//	var id string
//	err := p.database.QueryRow(sqlStatement,
//		floor.ID,
//		floor.Name,
//	).Scan(&id)
//	if err != nil {
//		return err
//	}
//	if id != floor.ID {
//		return CreateError
//	}
//	return nil
//}

//func (p PostgresDBStore) UpdateFloor(id string, floor *model.Floor) error {
//	sqlStatement :=
//		`UPDATE floors
//				SET name = $2,
//				WHERE id = $1
//				RETURNING id;`
//	var _id string
//	err := p.database.QueryRow(sqlStatement,
//		id,
//		floor.Name,
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
//func (p PostgresDBStore) RemoveFloor(id string) error {
//	sqlStatement :=
//		`DELETE FROM floors
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

func (p PostgresDBStore) queryMultipleFloors(sqlStatement string, args ...interface{}) ([]*model.Floor, error) {
	rows, err := p.database.Query(sqlStatement, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	floors := make([]*model.Floor, 0)
	for rows.Next() {
		var floor model.Floor
		err := rows.Scan(
			&floor.ID,
			&floor.Name,
		)
		if err != nil {
			// dont cause panic here, log it
			log.Printf("PostgresDBStore.queryMultipleOfferings: %v, sqlStatement: %s\n", err, sqlStatement)
		}
		floors = append(floors, &floor)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return floors, nil
}
