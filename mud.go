package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

// A MUD represents a MUD server.
type MUD struct {
	Clients map[net.Conn]Client
	ClientsMutex sync.RWMutex
	DB *DB
	Disconnects chan net.Conn
	Errors chan string
	Messages chan Message
	Server net.Listener
}

// Message represents something sent by a Client.
type Message struct {
	Conn net.Conn
	Message string
}

// BroadcastAll sends text to all clients.
func (m *MUD) BroadcastAll(s string) {
	// Convert message to bytes, to avoid repetitive conversion.
	b := []byte(s+"\n")

	m.ClientsMutex.RLock()
	defer m.ClientsMutex.RUnlock()
	for _, c := range m.Clients {
		c.Queue <- b
	}
}

// Handle handles a connection.
func (m *MUD) Handle(c net.Conn) {
	m.ClientsMutex.Lock()
	defer m.ClientsMutex.Unlock()

	// Spawn a new client.
	cl := Client{Conn: c}
	// Add it to our list.
	m.Clients[c] = cl
	// Send greeting.
	cl.Queue <- []byte(GREETING)
	// Serve it.
	cl.Serve(m.Disconnects)
	// Listen to it.
	cl.Listen(m.Messages, m.Disconnects)
}

// Listen opens and serves a listener.
func (m *MUD) Listen(port string) (err error) {
	m.Server, err = net.Listen("tcp", port)
	if err != nil {
		return
	}

	go func() {
		for {
			// Accept a new connection
			conn, err := m.Server.Accept()
			if err != nil {
				m.Errors <- fmt.Sprintf("error accepting connection, %e", err)
				continue
			}

			// Concurrently handle it
			go m.Handle(conn)
		}
	}()

	return nil
}

// Shutdown shuts the MUD down.
func (m *MUD) Shutdown() {
	for conn, client := range m.Clients {
		client.Queue <- []byte("[INFO] The MUD is shutting down. Goodbye.")
		conn.Close()
	}
}

// Start starts the MUD up.
func (m *MUD) Start(port string) (err error) {
	log.Println("[INFO] Lighthouse is starting...")

	// Perform initialization.
	m.Clients = make(map[net.Conn]Client)
	m.Disconnects = make(chan net.Conn)
	m.Errors = make(chan string)
	m.Messages = make(chan Message)

	// Open a port and listen for new connections.
	if err = m.Listen(port); err != nil {
		return
	}

	log.Println("[INFO] Lighthouse is operational.")
	// Enter the main loop.
	for {
		if err = m.Tick(); err != nil {
			return
		}
	}
}

// Tick does one round of processing.
func (m *MUD) Tick() (err error) {
	m.ClientsMutex.RLock()
	defer m.ClientsMutex.RUnlock()

	select {
	// Handle disconnects.
	case dc := <-m.Disconnects:
		cl, ok := m.Clients[dc]
		// The disconnected Conn should always be in the map, but always ensure, to avoid UB.
		if ok {
			m.ClientsMutex.Lock()
			delete(m.Clients, dc)
			m.ClientsMutex.Unlock()
			fmt.Println("left")
			m.BroadcastAll(fmt.Sprintf("%s has left the lighthouse.", cl.Name))
		}
	// Handle messages
	case msg := <-m.Messages:
		cl := m.Clients[msg.Conn]
		cl.Queue <- []byte("YOu say,"+msg.Message)
		// There's a message.
		m.BroadcastAll(msg.Message)
	default:
	}

	return
}
