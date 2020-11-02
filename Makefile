build:
	go build -o ./bin .

run:
	sh bin/start.sh

test:
	go test ./...
