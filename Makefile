.PHONY: clean tidy build-alpine

clean:
	rm -rf chat-api/bin/*
	rm -rf counter-api/bin/*

tidy:
	cd chat-api && go mod tidy
	cd counter-api && go mod tidy

build-alpine:
	cd chat-api && GOOS=linux GOARCH=amd64 go build -o bin/ ./...
	cd counter-api && GOOS=linux GOARCH=amd64 go build -o bin/ ./...

