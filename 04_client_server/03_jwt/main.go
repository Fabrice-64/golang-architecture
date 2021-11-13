package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/signin", Signin)
	http.ListenAndServe(":8080", nil)
}

func Index(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		log.Fatalln(http.StatusInternalServerError)
	}
}

func Signup(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "signup.gohtml", nil)
	if err != nil {
		log.Fatalln(http.StatusInternalServerError)
	}
}

func Signin(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "signin.gohtml", nil)
	if err != nil {
		log.Fatalln(http.StatusInternalServerError)
	}
}
