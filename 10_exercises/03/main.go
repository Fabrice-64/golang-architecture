package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	First    string
	Email    string
	Password []byte
}

var dbUser = map[string]User{} // email and User details
var tpl *template.Template
var key := []byte("My Secret Key")
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
		u = User{
			First:    f,
			Email:    e,
			Password: hp,
		}
		if err != nil {
			log.Println("Error while hashing the password")
		}

	}
	dbUser[u.Email] = u
	token := createToken(u.Email)

	log.Println("Token: ", token)
	_, err := parseToken(token)
	if err != nil {
		log.Println("Error while parsing the token")
	}
	tpl.ExecuteTemplate(w, "register.gohtml", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		e := req.FormValue("email")
		p := req.FormValue("password")
		log.Println("Method Post is used")
		if _, ok := dbUser[e]; ok {
			log.Println("User exists")
		} else {
			tpl.ExecuteTemplate(w, "login.html", nil)
			return
		}
		u := dbUser[e]
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			io.WriteString(w, "Hash and Password are different")
			log.Println("Hash and Password are different")
			return
		} else {
			log.Println("Hash and Password are identical")
			io.WriteString(w, "Hash and Password are identical")
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
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(sid))
	// HMAC needs to be converted into a printable form
	signedMac := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signedMac + "|" + sid

}

func parseToken(signedString string) (string, error) {
	xs := strings.SplitN(signedString, "|", 2)
	if len(xs) != 2{
		return "", "Stop hacking: wrong number of items!"
	}
	b64 := xs[0]
	// Has been encoded in base64 in createToken
	xb, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("Could not parse the token : ", err)
	}
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(xs[1]))
	ok :=hmac.Equal(xb, mac.Sum(nil))
	if !ok{
		return "", fmt.Error("Could not parse: different signed sid and sid")
	}
	return xs[1], nil
}
