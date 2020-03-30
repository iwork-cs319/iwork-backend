package utils

import (
	"database/sql"
	"io/ioutil"
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
func TimeYesterday() (time.Time, time.Time, error) {
	vancouver, err := time.LoadLocation("America/Vancouver")
	if err != nil {
	return time.Time{}, time.Time{}, err
	}
	twentyfour := time.Now().In(vancouver).Add(-24*time.Hour) // 24 Hours Ago
	yesterday := time.Date(twentyfour.Year(), twentyfour.Month(), twentyfour.Day(), 0, 0, 0, 0, twentyfour.Location()) // 0.0.0
	yesterdayEnd := yesterday.Add(time.Hour * time.Duration(23) + time.Minute * time.Duration(59) + time.Second * time.Duration(59)) // Time start of day + 23:59:59
	return yesterday, yesterdayEnd, nil
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
