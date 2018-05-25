#!/usr/bin/make -f

PACKAGE_NAME := docpile
PACKAGE_PATH := bitbucket.org/jonathanoliver/$(PACKAGE_NAME)

build:
	go fmt ./... && go build ./...

run:
	go run *.go

freeze:
	glock save -n "$(PACKAGE_PATH)" > .dependencies

restore: bitbucket
	cat .dependencies 2> /dev/null | glock sync -n "$(PACKAGE_PATH)"
