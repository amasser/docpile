#!/usr/bin/make -f

build:
	go fmt ./... && go build ./...

run:
	go run *.go
