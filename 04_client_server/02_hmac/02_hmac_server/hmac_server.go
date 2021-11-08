package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net"
)

// both client and server MUST have the same secret key
// to authenticate

var secret = "GolangIsAwesome"

func randStr(strSize int, randType string) string {
	var dict string
	if randType == "alphanum" {
		dict = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}
	if randType == "number" {
		dict = "0123456789"
	}
	if randType == "alpha" {
		dict = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}
	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dict[v%byte(len(dict))]
	}
	return string(bytes)
}

func serverSideAuthenticate(clientConn net.Conn, secretKey string) {
	msg := randStr(16, "alphanum")
	_, err := clientConn.Write([]byte(msg))
	if err != nil {
		log.Println("connection error: ", err)
		clientConn.Close()
	}
	hasher := hmac.New(md5.New, []byte(secretKey))
	hasher.Write([]byte(msg))
	serverHmacDigest := hasher.Sum(nil)
	fmt.Println("Server: ", base64.StdEncoding.EncodeToString(serverHmacDigest))
	buffer := make([]byte, 4096)
	n, err := clientConn.Read(buffer)
	if err != nil || n == 0 {
		fmt.Println("Error while reading buffer: ", err)
		clientConn.Close()
		return
	}
	clientHMacDigest := buffer[:n]
	fmt.Println("Relation Server-Client authenticated: ", hmac.Equal(serverHmacDigest, clientHMacDigest))
}

func handleConnection(c net.Conn) {
	log.Printf("Client %v connected", c.RemoteAddr())
	serverSideAuthenticate(c, secret)
	log.Printf("Connection from %v closed", c.RemoteAddr())
}

func main() {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server Listening on Port 6000")
	for {

		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}
