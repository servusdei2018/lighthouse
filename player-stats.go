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

// Name returns a Player's name.
func (s iPlayerStats) Name() string {
	s.RLock()
	defer s.RUnlock()
	return s.name
}
