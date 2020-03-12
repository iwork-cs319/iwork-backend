package utils

import (
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

func ParseSqlStatements(fileName string) ([]string, error) {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	stmts := strings.Split(string(dat), ";")
	return stmts, nil
}
