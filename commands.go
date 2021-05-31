package main

import (
	"fmt"
)

// Initialize the command map.
var CommandMap = make(map[string]func([]string, Player, MUD))
var CommandList []string

func init() {
	// Register commands.
	CommandMap["help"] = CmdHelp
	CommandMap["quit"] = CmdQuit
	CommandMap["who"] = CmdWho
	// Build command list.
	CommandList = make([]string, len(CommandMap))
	i := 0
	for k, _ := range CommandMap {
		CommandList[i] = k
		i++
	}
}

// CmdHelp shows available commands.
func CmdHelp(cmd []string, p Player, m MUD) {
	p.Send([]byte(fmt.Sprintf("Commands:\n%v\n", CommandList)))
}

// CmdQuit disconnects a Player from the game.
func CmdQuit(cmd []string, p Player, m MUD) {
	m.RemovePlayer(p)
}

// CmdWho lists all players currently connected.
func CmdWho(cmd []string, p Player, m MUD) {
	players := m.Players()
	p.Send([]byte("Players Online\n--------------\n"))
	for _, u := range players {
		if name := u.Name(); name != "" {
			p.Send([]byte(fmt.Sprintf("\t%s\n", name)))
		}
	}
	if amnt := len(players); amnt > 1 {
		p.Send([]byte(fmt.Sprintf("There are %v players online.\n", amnt)))
	} else {
		p.Send([]byte("There is 1 player online.\n"))
	}
}
