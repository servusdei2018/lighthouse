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
	"bufio"
	"net"
	"sync"
)

// Player describes a Player.
type Player interface {
	PlayerStats
	// Send sends a message to a Player.
	Send(msg []byte)
	// Shutdown closes the Player's connection.
	Shutdown()
}

// iPlayer implements Player.
type iPlayer struct {
	sync.RWMutex
	iPlayerStats
	// Player's network connection.
	conn net.Conn
	// Messages waiting to be sent to a player.
	msgqueue chan []byte
}

// NewPlayer creates a new Player from a network connection,
// and serves it.
func NewPlayer(conn net.Conn, mud MUD) Player {
	// Create a new Player.
	p := &iPlayer{conn: conn, msgqueue: make(chan []byte)}
	// Serve them.
	p.Serve(mud)
	// Send welcome message.
	p.Send(WELCOME_MSG)
	// Send welcome prompt.
	p.Send(WELCOME_PROMPT)
	return p
}

// Send sends a message to a Player.
func (p *iPlayer) Send(msg []byte) {
	p.msgqueue <- msg
}

// Serve serves a Player.
func (p *iPlayer) Serve(mud MUD) {
	// Listen for input.
	go func() {
		reader := bufio.NewReader(p.conn)
		var str string; var err error
		for {
			// Read one line of input.
			str, err = reader.ReadString('\n')
			// Relay the input, along with any error we encountered, to the MUD for processing.
			mud.Process(NewMessage(err, str, p))
			// Exit if we encountered an error.
			if err != nil { return }
		}
	}()
	// Send output.
	go func() {
		var err error
		for {
			if _, err = p.conn.Write(<-p.msgqueue); err != nil {
				// Relay the error to the MUD.
				mud.Process(NewMessage(err, "", p))
				// Exit.
				return
			}
		}
	}()
}

// Shutdown shuts a Player down.
func (p *iPlayer) Shutdown() {
	go p.conn.Close()
}
