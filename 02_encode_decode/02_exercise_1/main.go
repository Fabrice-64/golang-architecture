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
	http.HandleFunc("/encode", encode)
	http.HandleFunc("/decode", decode)
	http.ListenAndServe(":8080", nil)
}

func encode(w http.ResponseWriter, req *http.Request) {
	p1 := person{
		First: "Anne-Elisabeth",
	}
	err := json.NewEncoder(w).Encode(p1)
	if err != nil {
		log.Println("Not Encoded Value Issue ", err)
	}
	// curl localhost:8080/encode
}

func decode(w http.ResponseWriter, req *http.Request) {
	var p1 person
	err := json.NewDecoder(req.Body).Decode(&p1)
	if err != nil {
		log.Println("Not Decoded because Bad Value ", err)
	}
	log.Println("This value was passed: ", p1)
	// curl -XGET -H "Content-type: application/json" -d '{"First": "Papounet"}' 'localhost:8080/decode'
}
