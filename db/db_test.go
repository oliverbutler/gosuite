package db

import (
	"fmt"
	"reflect"
	"testing"
)

func TestExecute(t *testing.T) {
	db := Connect()
	defer db.Close()

	res, err := ExecuteSQL(db, "SELECT * FROM posts")
	if err != nil {
		panic(err)
	}

	if len(res.Rows) == 0 {
		t.Errorf("No rows returned")
	}

	fmt.Println(res)
}

func TestGetTable(t *testing.T) {
	db := Connect()
	defer db.Close()

	res, err := GetTables(db)
	if err != nil {
		panic(err)
	}

	expected := []string{"authors", "comments", "posts"}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Expected %v, got %v", expected, res)
	}
}
