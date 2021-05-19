build:
	go build -ldflags="-s -w"

deploy:
	nohup ./bin/lighthouse &
