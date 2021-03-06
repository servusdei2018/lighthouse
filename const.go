package main

import "fmt"

const (
	VERSION = "0.1.0-beta0.1"
)

var (
	WELCOME_MSG = []byte(fmt.Sprintf("Lighthouse %v (lighthouse.vineyard.haus:1990)\nPlease note that we are under heavy development.\n", VERSION))
	WELCOME_PROMPT = []byte("Welcome to the Lighthouse! By what name shall you be known? ")
	UNKNOWN_CMD_MSG = []byte("Huh?\n")
)
