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
	c, err := req.Cookie("sessionID")
	if err != nil {
		c = &http.Cookie{
			Name:  "sessionID",
			Value: "",
		}
	}
	s, err := parseToken(c.Value)
	if err != nil {
		log.Println(err)
	}
	var e string
	if s != "" {
		e = sessions[s]
	}
	log.Println("Email de l'utilisateur: ", e)
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func register(w http.ResponseWriter, req *http.Request) {
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
		bs, _ := json.Marshal(dbUser)
		log.Println("Registered Users: ", string(bs))
		// create a UUID - call the function
		sUUID := createUUid()
		// Connect UUID and user
		sessions[sUUID] = u.Email
		ss, _ := json.Marshal(sessions)
		log.Println("Sessions Users: ", string(ss))

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
