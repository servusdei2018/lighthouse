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
	"sync"
)

// PlayerStats describes in-game data for a Player.
type PlayerStats interface {
	Name() string
}

// iPlayerStats implements PlayerStats.
type iPlayerStats struct {
	sync.RWMutex
	// Player name.
	name string
}

// NewPlayerStats creates a new PlayerStats.
func NewPlayerStats(name string) {
	return &iPlayerStats{name: name}
}

// Name returns a Player's name.
func (s *iPlayerStats) Name() string {
	s.RLock()
	defer s.RUnlock()
	return s.name
}
