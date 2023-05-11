.PHONY: proto-chat-api proto-counter-api proto-all

clean:
	rm -rf chat-api/kl-proto/*
	rm -rf counter-api/kl-proto/*

proto-chat-api:
	protoc --go_out=chat-api --go_opt=paths=source_relative --go-grpc_out=./chat-api --go-grpc_opt=paths=source_relative ./kl-proto/*.proto

proto-counter-api:
	protoc --go_out=counter-api --go_opt=paths=source_relative --go-grpc_out=./counter-api --go-grpc_opt=paths=source_relative ./kl-proto/*.proto

proto-all: proto-chat-api proto-counter-api
