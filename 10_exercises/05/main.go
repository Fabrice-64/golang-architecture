package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	First    string
	Email    string
	Password []byte
}

type CustomClaims struct {
	jwt.StandardClaims
	sID string
}

var secretKey = "MySecretKey"
var dbUsers = map[string]User{}
var dbSession = map[string]string{}
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
	c, err := req.Cookie("Cookie-05")
	if err != nil {
		c = &http.Cookie{
			Name:  "Cookie-05",
			Value: "",
		}
		tpl.ExecuteTemplate(w, "index.gohtml", "Unknown")
		return
	}
	log.Println("Value of Cookie: ", c.Value)
	sID, err := parseToken(c.Value)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	e := dbSession[sID]
	u := dbUsers[e]
	log.Println("User email is: ", u.Email)
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func register(w http.ResponseWriter, req *http.Request) {
	// Get a User
	u := User{}
	if req.Method == http.MethodPost {
		f := req.FormValue("first")
		e := req.FormValue("email")
		p := req.FormValue("password")
		if f == "" || e == "" || p == "" {
			io.WriteString(w, "Please fill the fields")
			return
		}
		hp, err := hashPassword(p)
		if err != nil {
			log.Println("Sorry, coul not hash the pwd: ", err)
		}
		if _, ok := dbUsers[e]; ok {
			io.WriteString(w, "This user already exists")
			return
		}
		u = User{
			First:    f,
			Email:    e,
			Password: hp,
		}
		dbUsers[e] = u
		// Check if user got recorded
		bs, _ := json.Marshal(dbUsers)
		log.Println("Register - dbuser is: ", string(bs))
		// Create a UUID
		sUUID := uuid.NewString()
		// record the user UUID in a session.
		dbSession[sUUID] = u.Email
		// Create a cookie.
		c := &http.Cookie{
			Name:  "Cookie-05",
			Value: "",
		}
		// add a token to the Cookie
		token, err := createToken(sUUID)
		if err != nil {
			msg := url.QueryEscape("Could not create a token")
			http.Redirect(w, req, "/?msg="+msg, http.StatusInternalServerError)
		}
		c.Value = token
		http.SetCookie(w, c)
		log.Println("Register - Cookie: ", c)
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
			io.WriteString(w, "Please fill the form")
			return
		}
		// check if User exists
		if _, ok := dbUsers[e]; ok {
			log.Println("Login Function - User exists")
		} else {
			io.WriteString(w, "This User does not seem to exist")
			return
		}
		u := dbUsers[e]
		// As the User exists, check the password
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			io.WriteString(w, "Error in the password, obviously")
			return
		}
		// As the User exists & pwd is ok: create session
		sUUID := uuid.NewString()
		token, err := createToken(sUUID)
		if err != nil {
			msg := url.QueryEscape("Could not create a token")
			http.Redirect(w, req, "/?msg="+msg, http.StatusInternalServerError)
			return
		}
		c := &http.Cookie{
			Name:  "Cookie-05",
			Value: token,
		}
		http.SetCookie(w, c)
		tpl.ExecuteTemplate(w, "index.gohtml", u.First)
		return
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

func createToken(sUUID string) (string, error) {
	cc := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		sID: sUUID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	st, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error while signing the jwt token: %w", err)
	}
	return st, nil
}

func parseToken(st string) (string, error) {
	token, err := jwt.ParseWithClaims(st, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("could not parse the token %w", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("token parsing sent back wrong result")
	}
	return token.Claims.(*CustomClaims).sID, nil
}
