package main

import (
	"log"
	"net"
)

// Identify a process : sudo lsof -i :6000
// Kill it : kill -9 <PID>

func handleConnection(c net.Conn) {
	log.Printf("Client Connected : %v", c.RemoteAddr())
	log.Printf("Connection Closed: %v", c.RemoteAddr())
}

func main() {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalln("No connection established: ", err)
	}
	log.Println("Connection established on port 6000")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)

	}

}
