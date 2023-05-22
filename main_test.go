package main

import (
	"database/sql"
	"testing"
)

func TestConnectingToDB(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://postgres:postgrespw@localhost:32768/PokerGame?sslmode=disable")
	if err != nil {
		t.Errorf("we can't connect to DB")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("we can't connect to DB")
	}
}
