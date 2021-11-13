package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("my-text.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		log.Fatalln("could not io copy: ", err)
	}
	fmt.Printf("Type before Sum(nil): %T\n", h)
	fmt.Println(h)
	xb := h.Sum(nil)
	fmt.Printf("Type after Sum(nil): %T\n", xb)
	fmt.Printf("%x\n", xb)

}
