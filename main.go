package main

import (
	"log"
	"net"
)

func main() {
	// Start a listener at the specified port
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	log.Println("Lighthouse is operational.")

	// Accept connections
	for {
		// Accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
			continue
		}

		// Concurrently handle it
		go Handle(conn)
	}
}

// Handle handles a connection.
func Handle(conn net.Conn) {
	// Send welcome message
	conn.Write([]byte("Welcome to the Lighthouse! We are under construction, please monitor our website (https://lighthouse.vineyard.haus for information!"))
}
