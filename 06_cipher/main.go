package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var str = "Bonjour c'est vraiment épatant : j'aime découvrir de nouvelles options !"

func main() {
	password := "japprendsGolang"
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("encryption failed: ", err)
	}
	bs = bs[:16]
	res, err := EncryptDecode(bs, str)
	if err != nil {
		log.Println("error with encryptdecode func: ", err)
	}
	fmt.Println(res)
	fmt.Println(string(res))
}

func EncryptDecode(key []byte, msg string) ([]byte, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt: %w", err)
	}
	iv := make([]byte, aes.BlockSize)

	s := cipher.NewCTR(b, iv)
	buff := &bytes.Buffer{}
	sw := cipher.StreamWriter{
		S: s,
		W: buff,
	}
	_, err = sw.Write([]byte(msg))
	if err != nil {
		return nil, fmt.Errorf("error while writing: %w", err)
	}
	return buff.Bytes(), nil
}
