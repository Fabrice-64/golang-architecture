package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var key = []byte{}

func main() {
	for i := 1; i <= 64; i++ {
		key = append(key, byte(i))
	}
	p := "12345678"
	bs, err := hashPassword(p)
	if err != nil {
		log.Panic(err)
	}
	err = comparePassword(p, bs)
	if err != nil {
		log.Fatalln("Not Logged In")
	}
	log.Println("Logged In")

}

func hashPassword(p string) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error while hashing the password : %w", err)
	}
	return bs, nil
}

func comparePassword(p string, hp []byte) error {
	err := bcrypt.CompareHashAndPassword(hp, []byte(p))
	if err != nil {
		fmt.Errorf("Error in hashing the password: %w", err)
	}
	return nil
}

func signMsg(msg []byte) ([]byte, error) {
	h := hmac.New(sha512.New, key)
	_, err := h.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("not succeeded in Hashing Value: %w", err)
	}
	signed := h.Sum(nil)
	return signed, nil
}

func checkSig(msg, sig []byte) (bool, error) {
	newSig, err := signMsg(msg)
	if err != nil {
		return false, fmt.Errorf("error while signing the message: %w", err)
	}
	same := hmac.Equal(newSig, sig)
	return same, nil
}
