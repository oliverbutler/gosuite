package db

import (
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {
	db := Connect()

	res, err := ExecuteSQL(db, "SELECT * FROM posts")
	if err != nil {
		panic(err)
	}

	if len(res.Rows) == 0 {
		t.Errorf("No rows returned")
	}

	fmt.Println(res)

	db.Close()
}
