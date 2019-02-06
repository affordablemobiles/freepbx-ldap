package main

import (
	"fmt"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB

	sqlserver string = "127.0.0.1:3306"
	sqluser   string = "root"
	sqlpass   string = ""
	sqldb     string = "asterisk"
)

func SQLConnect() (err error) {
	if strings.ContainsRune(sqlserver, ':') {
		dbConn, err = sql.Open("mysql", sqluser+":"+sqlpass+"@tcp("+sqlserver+")/"+sqldb)
	} else {
		dbConn, err = sql.Open("mysql", sqluser+":"+sqlpass+"@tcp("+sqlserver+":3306)/"+sqldb)
	}
	if err != nil {
		err = dbConn.Ping()
	}
	return
}

type PhonebookEntry struct {
	Name      string
	Extension string
}

func SQLSearch(sqlQuery string, sqlVals []interface{}) ([]*PhonebookEntry, error) {
	var (
		rows   *sql.Rows
		err    error
		result []*PhonebookEntry = []*PhonebookEntry{}
	)
	rows, err = dbConn.Query(sqlQuery, sqlVals...)
	if err != nil {
		return nil, fmt.Errorf("Database Error: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name      string
			extension string
		)
		err := rows.Scan(&name, &extension)
		if err != nil {
			return nil, fmt.Errorf("Database Error: %s", err)
		}
		result = append(result, &PhonebookEntry{
			Name:      name,
			Extension: extension,
		})
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Database Error: %s", err)
	}

	return result, nil
}
