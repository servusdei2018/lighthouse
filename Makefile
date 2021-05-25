build:
	go build -ldflags="-s -w" -o ./bin/lighthouse

race:
	go build -ldflags="-s -w" --race -o ./bin/lighthouse

deploy:
	nohup ./bin/lighthouse &

local:
	./bin/lighthouse
