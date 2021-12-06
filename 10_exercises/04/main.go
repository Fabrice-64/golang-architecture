package main

import (
	"net/http"
	"text/template"
)

type User struct {
	First    string
	Email    string
	Password []byte
}

var dbUser = map[string]User{} //user email and User details
var secretKey = "This is a secret key"
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
	tpl.ExecuteTemplate(w, "register.gohtml", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}
