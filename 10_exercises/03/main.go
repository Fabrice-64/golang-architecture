package main

import (
	"html/template"
	"log"
	"net/http"
)

type User struct {
	First    string
	Email    string
	Password string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func register(w http.ResponseWriter, req *http.Request) {
	u := User{}
	if req.Method == http.MethodPost {
		u.First = req.FormValue("first")
		u.Email = req.FormValue("email")
		u.Password = req.FormValue("password")
	}
	log.Println("User Values: ", u)

	tpl.ExecuteTemplate(w, "register.gohtml", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}
