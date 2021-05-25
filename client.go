package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
)

type Client struct {
	sync.RWMutex
	Conn net.Conn
	Queue chan []byte
	// Playername
	Name string
}

// Listen listens to all messages and forwards them to appropriate channels for processing.
func (c *Client) Listen(msgs chan Message, dc chan net.Conn) {
	go func() {
		reader := bufio.NewReader(c.Conn)

		var msg string
		var err error
		for {
			// Read one line of input.
			msg, err = reader.ReadString('\n')
			if err != nil {
				// Give up on a bad connection.
				dc <- c.Conn
				c.Conn.Close()
				return
			}

			// Remove invalid characters.
			msg = strings.ReplaceAll(msg, "\r", "")
			// Add it to the message queue.
			msgs <- Message{Conn: c.Conn, Message: msg}
		}
	}()
}

// Serve sends messages to a client.
func (c *Client) Serve(dc chan net.Conn) {
	c.Queue = make(chan []byte)
	go func() {
		var msg []byte
		for {
			msg = <-c.Queue
			if c.send(msg, dc) != nil { return }
		}
	}()
}


// send sends a message to the client.
func (c *Client) send(b []byte, dc chan net.Conn) (err error) {
	_, err = c.Conn.Write(b)
	// If there s an error, register this client as disconnected.
	if err != nil {
		dc <- c.Conn
		c.Conn.Close()
	}
	return
}
