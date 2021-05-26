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
	"strings"
	"sync"
)

// A message from a Player to the MUD.
type Message interface {
	Error() error
	Message() string
	Player() Player
}

// Implementation of Message.
type iMessage struct {
	sync.RWMutex
	err error
	msg string
	p Player
}

// NewMessage creates a new Message.
func NewMessage(err error, msg string, p Player) Message {
	m := &iMessage{err: err, msg: msg, p: p}
	m.filter()
	return m
}

// Error returns a Message's error.
func (m *iMessage) Error() error {
	m.RLock()
	defer m.RUnlock()
	return m.err
}

// Message returns a Message's Message.
func (m *iMessage) Message() string {
	m.RLock()
	defer m.RUnlock()
	return m.msg
}

// Player returns a Message's Player.
func (m *iMessage) Player() Player {
	m.RLock()
	defer m.RUnlock()
	return m.p
}

// filter removes \r and \n from message text.
func (m *iMessage) filter() {
	m.Lock()
	defer m.Unlock()
	m.msg = strings.ReplaceAll(m.msg, "\r", "")
	m.msg = strings.ReplaceAll(m.msg, "\n", "")
}
