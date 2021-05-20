build:
	go build -ldflags="-s -w" -o ./bin/lighthouse

deploy:
	nohup ./bin/lighthouse &
	
local:
	./bin/lighthouse
