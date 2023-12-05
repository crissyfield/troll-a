NAME=troll-a
VERSION=`git describe --tag --always`
FLAGS=-ldflags="-s -w -X 'main.Version=${VERSION}'" -trimpath

COL_RESET=`tput sgr0`
COL_GREEN_BRIGHT=`tput setaf 10`
COL_YELLOW_BRIGHT=`tput setaf 11`

.PHONY: build clean

build: build-setup build-linux build-darwin

build-setup: clean
	@mkdir -p dist
	@go generate ./...

build-linux: build-linux-amd64 build-linux-arm64 build-linux-armv7

build-linux-amd64:
	@echo "Building ${COL_GREEN_BRIGHT}Linux-x86_64${COL_RESET}..."
	@GOOS=linux GOARCH=amd64 go build ${FLAGS} -o "dist/${NAME}-Linux-x86_64" .
	@md5 -r "dist/${NAME}-Linux-x86_64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Linux-x86_64" >> "dist/${NAME}-sha256sum.txt"

build-linux-arm64:
	@echo "Building ${COL_GREEN_BRIGHT}Linux-aarch64${COL_RESET}..."
	@GOOS=linux GOARCH=arm64 go build ${FLAGS} -o "dist/${NAME}-Linux-aarch64" .
	@md5 -r "dist/${NAME}-Linux-aarch64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Linux-aarch64" >> "dist/${NAME}-sha256sum.txt"

build-linux-armv7:
	@echo "Building ${COL_GREEN_BRIGHT}Linux-armv7l${COL_RESET}..."
	@GOOS=linux GOARCH=arm GOARM=7 go build ${FLAGS} -o "dist/${NAME}-Linux-armv7l" .
	@md5 -r "dist/${NAME}-Linux-armv7l" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Linux-armv7l" >> "dist/${NAME}-sha256sum.txt"

build-darwin: build-darwin-amd64 build-darwin-arm64

build-darwin-amd64:
	@echo "Building ${COL_GREEN_BRIGHT}Darwin-x86_64${COL_RESET}..."
	@GOOS=darwin GOARCH=amd64 go build ${FLAGS} -o "dist/${NAME}-Darwin-x86_64" .
	@md5 -r "dist/${NAME}-Darwin-x86_64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Darwin-x86_64" >> "dist/${NAME}-sha256sum.txt"

build-darwin-arm64:
	@echo "Building ${COL_GREEN_BRIGHT}Darwin-arm64${COL_RESET}..."
	@GOOS=darwin GOARCH=arm64 go build ${FLAGS} -o "dist/${NAME}-Darwin-arm64" .
	@md5 -r "dist/${NAME}-Darwin-arm64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Darwin-arm64" >> "dist/${NAME}-sha256sum.txt"

clean:
	@echo "${COL_YELLOW_BRIGHT}Cleaning dist${COL_RESET}..."
	@rm -rf dist
