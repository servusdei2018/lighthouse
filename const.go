package main

import "fmt"

const (
	VERSION = "0.1.0"
)

var (
	WELCOME_MSG = []byte(fmt.Sprintf("Lighthouse %v\nWelcome to the Lighthouse! By what name shall you be known? ", VERSION))
)
