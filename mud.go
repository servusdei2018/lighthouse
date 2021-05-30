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
	"fmt"
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
	// RemovePlayer removes a Player from the MUD's Player list.
	RemovePlayer(Player)
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

// RemovePlayer removes a Player from the MUD's player list.
func (m *iMUD) RemovePlayer(p Player) {
	// Shutdown the underlying connection.
	p.Shutdown()

	// Send a disconnect message.
	if p.Name() != "" {
		go m.BroadcastAll([]byte(fmt.Sprintf("%s has left the lighthouse.\n", p.Name())))
	}

	// Delete the Player from the player list.
	players := m.Players()
	for k, v := range players {
		if v == p {
			m.Lock()
			// Remove the Player by plucking him out of the list.
			m.players = append(m.players[:k], m.players[k+1:]...)
			m.Unlock()
			return
		}
	}
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
			go m.RemovePlayer(msg.Player())
			continue
		}

		if msg.Player().Name() != "" {
			go m.BroadcastAll([]byte(fmt.Sprintf("%v chats, \"%s\"\n", msg.Player().Name(), msg.Message())))
		} else {
			if msg.Message() != "" {
				msg.Player().SetName(msg.Message())
			} else {
				msg.Player().Send([]byte("You must choose a valid name!\n\n"))
				msg.Player().Send(WELCOME_MSG)
			}
		}
	}
}
