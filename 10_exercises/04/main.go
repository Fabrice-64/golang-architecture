package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	if req.Method == http.MethodPost {
		e := req.FormValue("email")
		p := req.FormValue("password")
		if e == "" || p == "" {
			io.WriteString(w, "please type a username/password")
		}
		if _, ok := dbUser[e]; ok {
			log.Println("User exists: ", e)
		} else {
			io.WriteString(w, "Error in username")
			return
		}
		u := dbUser[e]
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			io.WriteString(w, "Password and Hash do not match")
			return
		} else {
			tpl.ExecuteTemplate(w, "index.gohtml", nil)
			return
		}
	}
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

func parseToken(signedString string) (string, error) {
	xs := strings.SplitN(signedString, "|", 2)
	if len(xs) < 2 {
		return "", fmt.Errorf("Stop hacking my script")
	}
	b64 := xs[0]
	xb, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("Unable to parse the token: %s", err)
	}
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(xs[2]))
	ok := hmac.Equal(xb, mac.Sum(nil))
	if !ok {
		return "", fmt.Errorf("Could not parse different signed string and sid")
	}
}
