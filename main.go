package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var db DB
	if err := db.Import("db.json"); err != nil {
		panic(err)
	}

	var mud MUD
	mud.Start(":1990")
}
