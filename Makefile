.DEFAULT_GOAL := help

build:
	go build -ldflags="-s -w"

help:
	@echo "Makefile targets:"
	@echo ""
	@echo "build - build lighthouse for production"
	@echo "race - build lighthouse with race detection"
	@echo "test - run tests"
	@echo ""

race:
	go build -ldflags="-s -w" --race

test:
	go test
