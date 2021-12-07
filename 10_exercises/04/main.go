package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
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
	u := User{}
	if req.Method == http.MethodPost {
		f := req.FormValue("first")
		e := req.FormValue("email")
		p := req.FormValue("password")
		hp, err := hashPassword(p)
		if err != nil {
			log.Println("Error while hashing the Pwd")
		}
		u = User{
			First:    f,
			Email:    e,
			Password: hp,
		}
	}
	dbUser[u.Email] = u
	log.Println(dbUser[u.Email])
	token := createToken(u.Email)
	log.Println("Token: ", token)
	tpl.ExecuteTemplate(w, "register.gohtml", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func hashPassword(p string) ([]byte, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hp, nil
}

func createToken(sid string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(sid))
	signedMac := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signedMac + "|" + sid
}
