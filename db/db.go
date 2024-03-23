package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	config := mysql.Config{
		User:   "user",
		Passwd: "password",
		Addr:   "127.0.0.1:3306",
		DBName: "test",
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected!")

	return db
}

type executeResult struct {
	rows []map[string]interface{}
}

func ExecuteSQL(db *sql.DB, sql string) (executeResult, error) {
	rows, err := db.Query(sql)
	if err != nil {
		return executeResult{}, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return executeResult{}, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columnsData := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columnsData {
			columnPointers[i] = &columnsData[i]
		}

		// Scan the data into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return executeResult{}, err
		}

		// Create our map, and populate it with our data.
		rowData := make(map[string]interface{})
		for i, colName := range columns {
			val := columnPointers[i].(*interface{})
			rowData[colName] = *val
		}

		results = append(results, rowData)
	}
	if err := rows.Err(); err != nil {
		return executeResult{}, err
	}

	return executeResult{rows: results}, nil
}
