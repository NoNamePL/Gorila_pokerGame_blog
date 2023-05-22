package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	id   int    `json:"id"`
	Name string `json:"name"`
	Age  uint16 `json:"age"`
}

type Articles struct {
	Id       uint16 `json:"id"`
	Title    string ` json:"title"`
	Anons    string ` json:"anons"`
	FullText string ` json:"fullText"`
}

var posts = []Articles{}
var showPost = Articles{}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("postgres", "postgres://postgres:postgrespw@localhost:32768/PokerGame?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	query, err := db.Query(`SELECT * FROM "articles"`)
	if err != nil {
		panic(err)
	}

	posts = []Articles{}

	for query.Next() {
		var post Articles
		err = query.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	tmpl.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "create", nil)
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tmpl, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("postgres", "postgres://postgres:postgrespw@localhost:32768/PokerGame?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query, err := db.Query(fmt.Sprintf(`SELECT * FROM "articles" WHERE "id" = '%s'`,vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Articles{}
	// Выборка данных
	for query.Next() {
		var post Articles
		err = query.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		showPost = post
	}

	tmpl.ExecuteTemplate(w, "show", showPost)

}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	if title == "" || anons == "" || fullText == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		db, err := sql.Open("postgres", "postgres://postgres:postgrespw@localhost:32768/PokerGame?sslmode=disable")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			panic(err)
		}

		insertSmt := `insert into "articles"("title","anons","full_text") VALUES($1,$2,$3)`
		_, err = db.Exec(insertSmt, title, anons, fullText)
		if err != nil {
			panic(err)
		}
		http.Redirect(w, r, "../", http.StatusSeeOther)
	}
}
func handleRequest() {

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/saveArticle", saveArticle).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequest()

	/*
		//_, err = db.Exec(fmt.Sprintf("INSERT INTO 'articles' VALUES (DEFAULT,'%s','%s','%s')", title, anons, fullText))
		insert, err := db.Query(fmt.Sprintf(`INSERT INTO articles VALUES (DEFAULT,'%s','%s','%s')`, title, anons, fullText))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
	*/

	// выборка данных
	/*
		res, err := db.Query(`SELECT * FROM "users"`)
		if err != nil {
			panic(err)
		}

		for res.Next() {
			var user User
			// проверка существует ли определенное значение
			err = res.Scan(&user.id, &user.Name, &user.Age)
			if err != nil {
				panic(err)
			}

			fmt.Println(fmt.Sprintf("User id: %d, User name: %s, User age: %d", user.id, user.Name, user.Age))
		}

		fmt.Println("Succefully connect to db")
	*/
}
