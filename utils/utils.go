package utils

import (
	"database/sql"
	"io/ioutil"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const descEnd = "\n For additional questions contact icbc at cs319.icbc@gmail.com"
const icsTemplate = `
BEGIN:VCALENDAR
PRODID:-//Microsoft Corporation//Outlook 12.0 MIMEDIR//EN
VERSION:2.0
METHOD:PUBLISH
X-MS-OLK-FORCEINSPECTOROPEN:TRUE
BEGIN:VTIMEZONE
TZID:Central Pacific Standard Time
BEGIN:STANDARD
DTSTART:19500402T020000
TZOFFSETFROM:-0700
TZOFFSETTO:-0800
RRULE:FREQ=YEARLY;BYMINUTE=0;BYHOUR=2;BYDAY=1SU;BYMONTH=4
END:STANDARD
BEGIN:DAYLIGHT
DTSTART:19501001T020000
TZOFFSETFROM:-0800
TZOFFSETTO:-0700
RRULE:FREQ=YEARLY;BYMINUTE=0;BYHOUR=2;BYDAY=1SU;BYMONTH=10
END:DAYLIGHT
END:VTIMEZONE
BEGIN:VEVENT
CLASS:PUBLIC
DESCRIPTION:%s
DTSTART;TZID="Central Pacific Standard Time":%s
DTEND;TZID="Central Pacific Standard Time":%s
LOCATION:%s
PRIORITY:5
SEQUENCE:0
SUMMARY;LANGUAGE=en-us:%s
TRANSP:OPAQUE
UID:040000008200E00074C5B7101A82E008000000008062306C6261CA01000000000000000
X-MICROSOFT-CDO-BUSYSTATUS:FREE
X-MICROSOFT-CDO-IMPORTANCE:1
X-MICROSOFT-DISALLOW-COUNTER:FALSE
X-MS-OLK-ALLOWEXTERNCHECK:TRUE
X-MS-OLK-AUTOFILLLOCATION:FALSE
X-MS-OLK-CONFTYPE:0
BEGIN:VALARM
TRIGGER:-PT1440M
ACTION:DISPLAY
DESCRIPTION:Reminder
END:VALARM
END:VEVENT
END:VCALENDAR
`

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
