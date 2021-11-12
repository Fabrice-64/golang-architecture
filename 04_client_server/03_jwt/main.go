package main

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var router *mux.Router

func CreateRouter() {
	router = mux.NewRouter()
}

func InitializeRoutes() {
	router.HandleFunc("/signup", SignUp).Methods("POST")
	router.HandleFunc("/signin", Signin).Methods("POST")

}

func main() {
	CreateRouter()
	InitializeRoutes()
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}
