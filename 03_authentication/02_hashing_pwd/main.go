package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	p := "1234678"

	hp, err := hashPassword(p)
	if err != nil {
		log.Panic(err)
	}
	err = comparePassword(p, hp)
	if err != nil {
		log.Fatalln("not Logged In")
	}
	log.Println("Logged In")
}

func hashPassword(p string) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error Using bcrypt package : %w", err)
	}
	return bs, nil
}

func comparePassword(p string, hp []byte) error {
	err := bcrypt.CompareHashAndPassword(hp, []byte(p))
	if err != nil {
		return fmt.Errorf("error between password and hashed password : %w", err)
	}
	return nil
}
