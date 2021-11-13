package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

var str = "Bonjour c'est vraiment épatant : j'aime découvrir de nouvelles options !"

func main() {
	encoded := encode(str)
	fmt.Println(encoded)
	decoded, err := decode(encoded)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(decoded)
}

func encode(str string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	return encoded
}

func decode(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("could not decode %w", err)
	}
	return string(decoded), err
}
