package postgres

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go-api/db"
	"go-api/utils"
	"os"
	"testing"
	"time"
)

const (
	MainFloorId = "a709a01e-e9e5-3139-8f0d-fba4c0a2187f"
	Workspace1  = "a12411d3-d281-3735-b000-bf94b094d2af"
	Workspace2  = "f2188d5e-509f-3086-a031-d86da93a55c4"
	Workspace3  = "62eb266e-740d-3973-9e50-77b0287e3026"
	Workspace4  = "5e56de3d-2323-372d-897f-23d6037c8581"
	Workspace5  = "aad40cbb-4baf-3931-a5d2-6f98b414182a"
	Workspace6  = "bb15369d-e6e0-33b8-8b97-1779f8865890"
	Workspace7  = "3361d373-781a-34d7-bbb8-c7d562a0cf51"
)

type AvailabilityTestSuite struct {
	suite.Suite
	store *db.DataStore
}

func (suite *AvailabilityTestSuite) SetupSuite() {
	dbUrl := os.Getenv("TEST_DB_URL")
	fixtures := []string{
		"../../resources/tables.sql",
		"../../test-fixtures/floors.sql",
		"../../test-fixtures/users.sql",
		"../../test-fixtures/workspaces.sql",
		"../../test-fixtures/book_offer.sql",
	}
	if err := utils.RunFixturesOnDB(dbUrl, fixtures); err != nil {
		suite.FailNow("failed to re-seed test db")
	}
	store, err := NewPostgresDataStore(dbUrl)
	if err != nil {
		suite.FailNow("failed to connect to DB" + err.Error())
	}
	suite.store = store
}

func TestAvailability(t *testing.T) {
	suite.Run(t, new(AvailabilityTestSuite))
}

func date(str string) time.Time {
	parse, _ := time.Parse(time.RFC3339, str)
	return parse
}

func (suite *AvailabilityTestSuite) TestFindAvailability1() {
	store := suite.store
	t := suite.T()
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-13T00:00:00Z"),
		date("2019-01-16T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}
	assert.Equal(t, 2, len(ids))
	assert.Contains(t, ids, Workspace2)
	assert.Contains(t, ids, Workspace7)
}

func (suite *AvailabilityTestSuite) TestFindAvailability2() {
	store := suite.store
	t := suite.T()
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-13T00:00:00Z"),
		date("2019-01-17T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}
	assert.Equal(t, 1, len(ids))
	assert.Contains(t, ids, Workspace7)
}

func (suite *AvailabilityTestSuite) TestFindAvailability3() {
	store := suite.store
	t := suite.T()
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-15T00:00:00Z"),
		date("2019-01-16T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}

	assert.Equal(t, 3, len(ids))
	assert.Contains(t, ids, Workspace2)
	assert.Contains(t, ids, Workspace3)
	assert.Contains(t, ids, Workspace7)
}

func (suite *AvailabilityTestSuite) TestFindAvailability4() {
	store := suite.store
	t := suite.T()
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-19T00:00:00Z"),
		date("2019-01-21T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}

	assert.Equal(t, 3, len(ids))
	assert.Contains(t, ids, Workspace1)
	assert.Contains(t, ids, Workspace2)
	assert.Contains(t, ids, Workspace6)
}

func (suite *AvailabilityTestSuite) TestFindAvailability5() {
	store := suite.store
	t := suite.T()
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-24T00:00:00Z"),
		date("2019-01-27T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}

	assert.Equal(t, 2, len(ids))
	assert.Contains(t, ids, Workspace2)
	assert.Contains(t, ids, Workspace3)
}

func (suite *AvailabilityTestSuite) TestFindAvailability6() {
	store := suite.store
	t := suite.T()
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-22T00:00:00Z"),
		date("2019-01-24T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}

	assert.Equal(t, 3, len(ids))
	assert.Contains(t, ids, Workspace3)
	assert.Contains(t, ids, Workspace5)
	assert.Contains(t, ids, Workspace7)
}

func (suite *AvailabilityTestSuite) TestFindAvailability8() {
	store := suite.store
	t := suite.T()
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-17T00:00:00Z"),
		date("2019-01-22T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}

	assert.Equal(t, 2, len(ids))
	assert.Contains(t, ids, Workspace1)
	assert.Contains(t, ids, Workspace6)
}

func (suite *AvailabilityTestSuite) TestFindAvailability9() {
	store := suite.store
	t := suite.T()
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-17T00:00:00Z"),
		date("2019-01-24T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}
	assert.Equal(t, 0, len(ids))
}
