.DEFAULT_GOAL := goapp

.PHONY: all
all: clean goapp

.PHONY: server
server:
	mkdir -p bin
	go build -o bin/server cmd/server/main.go

.PHONY: client
client:
	mkdir -p bin
	go build -o bin/client cmd/client/main.go

.PHONY: goapp
goapp: server client

.PHONY: clean
clean:
	go clean
	rm -f bin/*
