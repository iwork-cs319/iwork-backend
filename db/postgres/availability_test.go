package postgres

import (
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

func date(str string) time.Time {
	parse, _ := time.Parse(time.RFC3339, str)
	return parse
}

func testArrayEquality(t *testing.T, arr1, arr2 []string) {
	if len(arr1) != len(arr2) {
		t.Fatalf("arrays of different length, arr1=%s, arr2=%s", arr1, arr2)
	}
	for _, e1 := range arr1 {
		found := false
		for _, e2 := range arr2 {
			if e1 == e2 {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("%s not found in arr1=%s, arr2=%s", e1, arr1, arr2)
		}
	}
}

func TestFindAvailability1(t *testing.T) {
	store := CreateTestDBConn(t)
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-13T00:00:00Z"),
		date("2019-01-16T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}
	testArrayEquality(t, ids, []string{Workspace2, Workspace7})
}

func TestFindAvailability2(t *testing.T) {
	store := CreateTestDBConn(t)
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-13T00:00:00Z"),
		date("2019-01-17T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}
	testArrayEquality(t, ids, []string{Workspace7})
}

func TestFindAvailability3(t *testing.T) {
	store := CreateTestDBConn(t)
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-15T00:00:00Z"),
		date("2019-01-16T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}
	testArrayEquality(t, ids, []string{Workspace2, Workspace3, Workspace7})
}

func TestFindAvailability4(t *testing.T) {
	store := CreateTestDBConn(t)
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-19T00:00:00Z"),
		date("2019-01-21T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}
	testArrayEquality(t, ids, []string{Workspace1, Workspace2, Workspace6})
}

func TestFindAvailability5(t *testing.T) {
	store := CreateTestDBConn(t)
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-24T00:00:00Z"),
		date("2019-01-27T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}
	testArrayEquality(t, ids, []string{Workspace2, Workspace3})
}

func TestFindAvailability6(t *testing.T) {
	store := CreateTestDBConn(t)
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-22T00:00:00Z"),
		date("2019-01-24T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}
	testArrayEquality(t, ids, []string{Workspace3, Workspace5, Workspace7})
}

func TestFindAvailability8(t *testing.T) {
	store := CreateTestDBConn(t)
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-17T00:00:00Z"),
		date("2019-01-22T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}
	testArrayEquality(t, ids, []string{Workspace1, Workspace6})
}

func TestFindAvailability9(t *testing.T) {
	store := CreateTestDBConn(t)
	ids, err := store.WorkspaceProvider.FindAvailability(
		MainFloorId,
		date("2019-01-17T00:00:00Z"),
		date("2019-01-24T00:00:00Z"),
	)
	if err != nil {
		t.Fatal("error getting ids", err)
	}
	testArrayEquality(t, ids, []string{})
}
