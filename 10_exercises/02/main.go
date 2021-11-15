package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	First    string
	Email    string
	Password string
}

var tpl *template.Template

var dbUser = map[string]User{} //user Email, user encrypted pwd

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
	if req.Method != http.MethodPost {
		http.Redirect(w, req, "/", http.StatusMethodNotAllowed)
	}
	first := req.FormValue("first")
	email := req.FormValue("email")
	password := req.FormValue("password")
	cryptedPwd, err := hashPassword(password)
	if err != nil {
		fmt.Println("error while using bcrypt: %w", err)
	}

	user := User{
		First:    first,
		Email:    email,
		Password: string(cryptedPwd),
	}
	dbUser[user.Email] = user

	log.Printf("First Name: %s, Email: %s, Password: %s", user.First, user.Email, user.Password)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func login(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func hashPassword(pwd string) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
