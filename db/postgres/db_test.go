package postgres

import (
	"go-api/db"
	"os"
	"testing"
)

func CreateTestDBConn(t *testing.T) *db.DataStore {
	dbUrl := os.Getenv("TEST_DB_URL")
	store, err := NewPostgresDataStore(dbUrl)
	if err != nil {
		t.Fatal("failed to connect to DB" + err.Error())
	}
	return store
}

func TestDBCreate(t *testing.T) {
	conn := CreateTestDBConn(t)
	conn.Close()
}
