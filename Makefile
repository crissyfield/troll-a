NAME=troll-a
VERSION=`git describe --tag --always`
FLAGS=-ldflags="-s -w -X 'main.Version=${VERSION}'" -trimpath -tags re2_cgo

.PHONY: build

build:
	@go generate ./...
	@go build ${FLAGS} -o "${NAME}" .
