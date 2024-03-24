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

	return db
}

type ExecuteResult struct {
	Rows    []map[string]interface{}
	Columns []string
}

func ExecuteSQL(db *sql.DB, sql string) (ExecuteResult, error) {
	rows, err := db.Query(sql)
	if err != nil {
		return ExecuteResult{}, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return ExecuteResult{}, err
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
			return ExecuteResult{}, err
		}

		// Create our map, and populate it with our data.
		rowData := make(map[string]interface{})
		for i, colName := range columns {
			val := columnPointers[i].(*interface{})

			// if int array make a string
			if *val == nil {
				rowData[colName] = nil
			} else if _, ok := (*val).([]byte); ok {
				rowData[colName] = string((*val).([]byte))
			} else {
				rowData[colName] = *val
			}
		}

		results = append(results, rowData)
	}
	if err := rows.Err(); err != nil {
		return ExecuteResult{}, err
	}

	return ExecuteResult{Rows: results, Columns: columns}, nil
}

func GetTableSchema(db *sql.DB, table string) (ExecuteResult, error) {
	res, err := ExecuteSQL(db, fmt.Sprintf("DESCRIBE %s", table))

	return res, err
}

func GetTables(db *sql.DB) ([]string, error) {
	tables := make([]string, 0)

	res, err := ExecuteSQL(db, "SHOW TABLES")
	if err != nil {
		return tables, err
	}

	for _, row := range res.Rows {
		for _, col := range res.Columns {
			tables = append(tables, row[col].(string))
		}
	}

	return tables, nil
}
