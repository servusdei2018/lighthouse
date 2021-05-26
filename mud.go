/*
Lighthouse
Copyright (C) 2021  Nathanael Bracy

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"net"
	"sync"
)

// MUD facilitates interaction with its underlying components.
type MUD interface {
	// BroadcastAll sends a message to all Players.
	BroadcastAll(msg []byte)
	// ListenAndServe starts the MUD on the specified port.
	ListenAndServe(port string) error
	// Players returns a slice of online Players.
	Players() []Player
	// Process receives a Message for processing.
	Process(msg Message)
}

// MUD implementation.
type iMUD struct {
	sync.RWMutex
	msgqueue chan Message
	// Listener.
	server net.Listener
	// Player store.
	players []Player
}

// NewMud creates a new MUD.
func NewMud() MUD {
	return &iMUD{msgqueue: make(chan Message)}
}

// BroadcastAll sends a message to all Players.
func (m *iMUD) BroadcastAll(msg []byte) {
	for _, p := range m.Players() {
		p.Send(msg)
	}
}

// Players returns a slice of players currently online.
func (m *iMUD) Players() []Player {
	m.RLock()
	defer m.RUnlock()
	return m.players
}

// Process receives a Messsage for processing.
func (m *iMUD) Process(msg Message) {
	go func() {
		m.msgqueue <- msg
	}()
}

// ListenAndServe starts a MUD up.
func (m *iMUD) ListenAndServe(port string) (err error) {
	// Start the listener.
	m.Lock()
	m.server, err = net.Listen("tcp", port)
	if err != nil {
		m.Unlock()
		return
	}
	m.Unlock()

	// Accept connections.
	go func() {
		var conn net.Conn
		for {
			// Accept a new connection.
			conn, err = m.server.Accept()
			if err != nil { return }
			// Add a new player to the player pool.
			m.Lock()
			m.players = append(m.players, NewPlayer(conn, m))
			m.Unlock()
		}
	}()

	// Process messages.
	var msg Message
	for {
		msg = <- m.msgqueue
		if msg.Error() != nil {
			msg.Player().Shutdown()
			//TODO: Remove this player from the player pool
			continue
		}
		m.BroadcastAll([]byte(msg.Message()+"\n"))
	}
}
