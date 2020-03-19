package utils

import (
	"database/sql"
	"io/ioutil"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func TimeStampToTime(timestamp string) (time.Time, error) {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	tm := time.Unix(i, 0).UTC()
	return tm, err
}

func CalendarGenerator(startTime time.Time, endTime time.Time, desc string, summary string, location string) {
	// Title Booking for Workspace + id
	// TODO: Build link
	fmt.Sprintf(icsTemplate, desc+descEnd, startTime, endTime, location, summary)
}

func RunFixturesOnDB(dbUrl string, fileNames []string) error {
	database, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	if err = database.Ping(); err != nil {
		return err
	}
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, file := range fileNames {
		if err = execStmts(tx, file); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func execStmts(tx *sql.Tx, fileName string) error {
	statements, err := parseSqlStmts(fileName)
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		_, err = tx.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseSqlStmts(fileName string) ([]string, error) {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	stmts := strings.Split(string(dat), ";")
	return stmts, nil
}
