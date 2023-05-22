package main

import (
	"database/sql"
	"net/http"
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

func TestConnectingServer(t *testing.T) {
	http.HandleFunc("/", index)
}

func TestCheckToGetAllArticles(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://postgres:postgrespw@localhost:32768/PokerGame?sslmode=disable")
	if err != nil {
		t.Errorf("we can't connect to DB")
	}
	defer db.Close()
	query, err := db.Query(`SELECT * FROM "articles"`)
	if err != nil {
		t.Errorf("we can't send quares to DB")
	}

	type Articles struct {
		Id       uint16 `json:"id"`
		Title    string ` json:"title"`
		Anons    string ` json:"anons"`
		FullText string ` json:"fullText"`
	}
	for query.Next() {
		// проверка существует ли определенное значение
		var post Articles
		err = query.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			t.Errorf("we can't get dates to DB")
		}

	}
}
