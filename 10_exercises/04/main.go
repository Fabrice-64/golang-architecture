package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	First    string
	Email    string
	Password []byte
}

var dbUser = map[string]User{}     //user email and User details
var sessions = map[string]string{} //session ID (UUID) and user ID (email)
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
	log.Println("We are currently on Index function !")
	c, err := req.Cookie("Cookie-04")
	if err != nil {
		c = &http.Cookie{
			Name:  "Cookie-04",
			Value: "",
		}
	}
	log.Println("Index - Cookie-04 Value: ", c)
	s, err := parseToken(c.Value)
	if err != nil {
		log.Println(err)
	}
	var e string
	if s != "" {
		e = sessions[s]
	}
	u := dbUser[e]
	log.Println("Prénom de l'utilisateur: ", u.First)
	tpl.ExecuteTemplate(w, "index.gohtml", u.First)
}

func register(w http.ResponseWriter, req *http.Request) {
	log.Println("We have just landed on register function !")
	u := User{}
	if req.Method == http.MethodPost {
		log.Println("Method is POST !")
		e := req.FormValue("email")
		f := req.FormValue("first")
		p := req.FormValue("password")
		hp, err := hashPassword(p)
		if err != nil {
			log.Println("could hash the pwd: ", err)
		}
		if _, ok := dbUser[e]; ok {
			io.WriteString(w, "User Already exists")
			return
		}
		u = User{
			First:    f,
			Email:    e,
			Password: hp,
		}
		dbUser[u.Email] = u
		// for check
		bs, _ := json.Marshal(dbUser)
		log.Println("Registered Users: ", string(bs))
		// create a UUID - call the function
		sUUID := createUUid()
		// Connect UUID and user
		sessions[sUUID] = u.Email
		// for check
		ss, _ := json.Marshal(sessions)
		log.Println("Sessions Users: ", string(ss))
		// Create Token for the cookie: a hash | the uuid
		token := createToken(sUUID)
		// Create cookie
		c := &http.Cookie{
			Name:  "Cookie-04",
			Value: token,
		}
		// Set cookie
		http.SetCookie(w, c)
		// Then redirect to "/"
		tpl.ExecuteTemplate(w, "index.gohtml", u.First)
		return

	}
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
			sUUID := createUUid()
			token := createToken(sUUID)
			c := &http.Cookie{
				Name:  "Cookie-04",
				Value: token,
			}
			http.SetCookie(w, c)
			tpl.ExecuteTemplate(w, "index.gohtml", u.First)
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

func createToken(suuid string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(suuid))
	signedMac := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signedMac + "|" + suuid
}

func parseToken(signedString string) (string, error) {
	xs := strings.SplitN(signedString, "|", 2)
	if len(xs) < 2 {
		return "", fmt.Errorf("stop hacking my script")
	}
	b64 := xs[0]
	xb, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("unable to parse the token: %s", err)
	}
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(xs[1]))
	ok := hmac.Equal(xb, mac.Sum(nil))
	if !ok {
		return "", fmt.Errorf("could not parse different signed string and sid")
	}
	return xs[1], nil
}

func createUUid() string {
	// by default it's a V4 uuid in the google package
	uuid := uuid.NewString()
	return uuid
}
