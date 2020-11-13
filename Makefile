build:
	go build .

run: build
	./tsm

test:
	go test ./...
