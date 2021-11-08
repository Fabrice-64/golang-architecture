package main

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"net"
)

var secret = "GolangIsAwesome"

func clientSideAuthenticate(serverConn net.Conn, secret string, message string) {
	hasher := hmac.New(md5.New, []byte(secret))
	hasher.Write([]byte(message))
	clientHMacDigest := hasher.Sum(nil)
	fmt.Println("HMAC Digest send to server: ", base64.StdEncoding.EncodeToString(clientHMacDigest))
	n, err := serverConn.Write(clientHMacDigest)
	if err != nil || n == 0 {
		serverConn.Close()
		return
	}
}

func handleConnection(c net.Conn) {
	buffer := make([]byte, 4096)
	for {
		n, err := c.Read(buffer)
		if err != nil || n == 0 {
			c.Close()
			break
		}
		msg := string(buffer[:n])
		fmt.Println("\nData received from Server: ", msg)
		clientSideAuthenticate(c, secret, msg)
	}
	fmt.Printf("Connection from : %v Closed", c.RemoteAddr())
}

func main() {
	hostName := "localhost"
	portNumber := "6000"
	for {
		dialConn, err := net.Dial("tcp", hostName+":"+portNumber)
		if err != nil {
			log.Println("Connection Error: ", err)
			return
		}
		log.Printf("Connection OK between distant maching this machine and %s", hostName)
		log.Printf("Name of distant machine: %s", dialConn.RemoteAddr().String())
		log.Printf("Local Adress : %s", dialConn.LocalAddr().String())
		go handleConnection(dialConn)
	}
}
