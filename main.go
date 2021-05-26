package main

import (
	"log"
	"math/rand"
	"time"
)

const (
	PORT = ":1990"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var lighthouse = NewMud()
	log.Fatal(lighthouse.ListenAndServe(PORT))
}
