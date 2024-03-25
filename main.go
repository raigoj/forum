package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var tpl *template.Template
var db *sql.DB
var authenticated = true

var inMemorySession *Session

func main() {
	inMemorySession = NewSession()
	var err error
	tpl, err = template.ParseGlob("templates/*")
	if err != nil {
		panic(err)
	}
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginHandler)
	http.HandleFunc("/posts", postsHandler)
	http.HandleFunc("/comments", commentsHandler)
	http.HandleFunc("/create", createHandler)
	// http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/testing", testHandler)
	fmt.Println("Server is listening at port 8181")
	http.Handle("/statics/",
		http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics"))),
	)
	http.Handle("/templates/",
		http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))),
	)
	http.ListenAndServe("localhost:8181", nil)

}
