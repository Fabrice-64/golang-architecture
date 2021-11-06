package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type person struct {
	First string
}

func main() {
	http.HandleFunc("/encode", foo)
	http.HandleFunc("/decode", bar)
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, req *http.Request) {
	p1 := person{
		First: "Jenny",
	}

	err := json.NewEncoder(w).Encode(p1)
	if err != nil {
		log.Println("Encoded Bad Data: ", err)
	}

}

func bar(w http.ResponseWriter, req *http.Request) {
	var p1 person
	err := json.NewDecoder(req.Body).Decode(&p1)
	if err != nil {
		log.Println("Decoded Bad Data: ", err)
	}
	log.Println("Decoded Value: ", p1)
}
