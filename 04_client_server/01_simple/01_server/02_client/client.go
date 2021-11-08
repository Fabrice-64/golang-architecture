package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	hostName := "localhost"
	portNum := "6000"

	conn, err := net.Dial("tcp", hostName+":"+portNum)
	if err != nil {
		log.Fatalln("Connection Failed")
		return
	}
	fmt.Printf("Connection established between : %v and this machine", hostName)
	fmt.Printf("Remote Address : %s \n", conn.RemoteAddr().String())
	fmt.Printf("Local Address : %s \n", conn.LocalAddr().String())

}
