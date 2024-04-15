package db

import (
	"database/sql"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-sql-driver/mysql"
	config "gosuite/services/config"
)

type Connection interface {
	Status() string
	GetConfig() *config.DatabaseConfig
	GetConnection() *sql.DB
}

type ConnectionSuccess struct {
	config *config.DatabaseConfig
	db     *sql.DB
}

func (cs ConnectionSuccess) Status() string {
	return "Connected"
}

func (cs ConnectionSuccess) GetConfig() *config.DatabaseConfig {
	return cs.config
}

func (cs ConnectionSuccess) GetConnection() *sql.DB {
	return cs.db
}

type ConnectionError struct {
	config *config.DatabaseConfig
	reason string
}

func (ce ConnectionError) Status() string {
	return ce.reason
}

func (ce ConnectionError) GetConfig() *config.DatabaseConfig {
	return ce.config
}

func (ce ConnectionError) GetConnection() *sql.DB {
	return nil
}

type ConnectionPending struct {
	Config *config.DatabaseConfig
}

func (cp ConnectionPending) Status() string {
	return "Connecting..."
}

func (cp ConnectionPending) GetConfig() *config.DatabaseConfig {
	return cp.Config
}

func (cp ConnectionPending) GetConnection() *sql.DB {
	return nil
}

func Connect(info *config.DatabaseConfig) *sql.DB {
	config := mysql.Config{
		User:                    info.User,
		Passwd:                  info.Password,
		Addr:                    fmt.Sprintf("%s:%d", info.Host, info.Port),
		DBName:                  info.Database,
		AllowNativePasswords:    true,
		AllowCleartextPasswords: true,
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

type errMsg error

func ExucuteSQLCmd(sql string, conn *sql.DB) tea.Cmd {
	return func() tea.Msg {
		result, err := ExecuteSQL(conn, sql)
		if err != nil {
			return errMsg(err)
		}
		return result
	}
}

type ExecuteResult struct {
	Query        string
	Rows         []map[string]interface{}
	Columns      []string
	Microseconds int64
}

func ExecuteSQL(db *sql.DB, sql string) (ExecuteResult, error) {
	start := time.Now()

	rows, err := db.Query(sql)
	if err != nil {
		return ExecuteResult{}, err
	}

	elapsed := time.Since(start)

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

	return ExecuteResult{
		Query:        sql,
		Rows:         results,
		Columns:      columns,
		Microseconds: elapsed.Microseconds(),
	}, nil
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
