package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type Client struct {
	Conn net.Conn
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

// Send sends a message to a client.
//
// Sending a string is slightly less efficient than sending bytes.
// When sending a message to a particular client, we usually use a string.
// However, when broadcasting to all players, the message is converted
// to bytes once (see mud.BroadcastAll), and this method is used to
// increase efficiency.
func (c *Client) Send(s string, dc chan net.Conn) {
	go func() {
		_, err := c.Conn.Write([]byte(s+"\n"))
		// If there s an error, register this client as disconnected.
		if err != nil {
			dc <- c.Conn
			c.Conn.Close()
		}
	}()
}

// SendBytes sends a message to a client.
//
// Sending bytes is slightly more efficient than sending a string,
// which must be converted to bytes before transfer. When sending
// a message to a particular client, we usually use a string. However,
// when broadcasting to all players, the message is converted to
// bytes once (see mud.BroadcastAll), and this method is used to
// increase efficiency.
func (c *Client) SendBytes(b []byte, dc chan net.Conn) {
	go func() {
		_, err := c.Conn.Write(b)
		// If there s an error, register this client as disconnected.
		if err != nil {
			dc <- c.Conn
			c.Conn.Close()
		}
	}()
}
