all:
	go build -ldflags="-s -w"

race:
	go build -ldflags="-s -w" --race

test:
	go test
