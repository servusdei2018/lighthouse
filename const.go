package main

import "fmt"

const (
	VERSION = "0.1.0-beta0.1"
)

var (
	WELCOME_MSG = []byte(fmt.Sprintf("Lighthouse %v (lighthouse.vineyard.haus:1990)\nPlease note that we are under heavy development.\nWelcome to the Lighthouse! By what name shall you be known? ", VERSION))
)
