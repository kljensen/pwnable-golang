package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB
var usersWhoLoggedIn []string

func init() {
	log.Println("woot")
	database, _ = sql.Open("sqlite3", "./pwnable-datastore.db")
}

func homepage(w http.ResponseWriter, _ *http.Request) {
	var loginTemplate = template.Must(template.New("homepage").Parse(`<html> <body> <h1>Hi, please log in</h1>
	<form action="/private" method="POST">
	<input name="username" placeholder="username">
	<input name="password" type="password" placeholder="password">
	<input type="submit">
	</form>
	<h2>Here's other people who logged in</h2>
	<ul>{{range .}}<li>{{.}}</li>{{end}}</ul>
	</body> </html>
	`))
	loginTemplate.Execute(w, usersWhoLoggedIn)
}

func privateArea(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		u := r.PostFormValue("username")
		p := r.PostFormValue("password")
		usersWhoLoggedIn = append(usersWhoLoggedIn, u)
		query := "SELECT id, username  FROM users WHERE username='" + u + "' AND password='" + p + "'"
		var userID int
		var username string
		result := database.QueryRow(query)
		result.Scan(&userID, &username)
		if username != "" {
			fmt.Fprint(w, "Welcome to the private area, "+u)
		} else {
			fmt.Fprint(w, "Forbidden! This is a private area and "+u+" is not allowed here.")
		}
	} else {
		fmt.Fprint(w, "Forbidden!")
	}
}

func main() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/private", privateArea)
	http.ListenAndServe(":8080", nil)
}
