NAME=troll-a
VERSION=`git describe --tag --always`
FLAGS=-ldflags="-s -w -X 'main.Version=${VERSION}'" -trimpath
FLAGS_RE2=-tags re2_cgo

COL_RESET=`tput sgr0`
COL_GREEN_BRIGHT=`tput setaf 10`
COL_YELLOW_BRIGHT=`tput setaf 11`

.PHONY: build-std
build-std: build-setup build-local-std

.PHONY: build-re2
build-re2: build-setup build-local-re2

.PHONY: build-dist
build-dist: build-setup build-mkdist build-dist-linux build-dist-darwin

.PHONY: build-setup
build-setup:
	@go generate ./...

.PHONY: build-local-std
build-local-std:
	@echo "Building ${COL_GREEN_BRIGHT}local with standard regular expressions${COL_RESET}..."
	@CGO_ENABLED=1 go build ${FLAGS} -o "${NAME}" .

.PHONY: build-local-re2
build-local-re2:
	@echo "Building ${COL_GREEN_BRIGHT}local with go-re2 regular expressions${COL_RESET}..."
	@CGO_ENABLED=1 go build ${FLAGS} ${FLAGS_CGO} -o "${NAME}" .

.PHONY: build-mkdist
build-mkdist: clean
	@mkdir -p dist

.PHONY: build-dist-linux
build-dist-linux: build-dist-linux-amd64 build-dist-linux-arm64 build-dist-linux-armv7

.PHONY: build-dist-linux-amd64
build-dist-linux-amd64:
	@echo "Building ${COL_GREEN_BRIGHT}Linux-x86_64${COL_RESET}..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${FLAGS} -o "dist/${NAME}-Linux-x86_64" .
	@md5 -r "dist/${NAME}-Linux-x86_64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Linux-x86_64" >> "dist/${NAME}-sha256sum.txt"

.PHONY: build-dist-linux-arm64
build-dist-linux-arm64:
	@echo "Building ${COL_GREEN_BRIGHT}Linux-aarch64${COL_RESET}..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ${FLAGS} -o "dist/${NAME}-Linux-aarch64" .
	@md5 -r "dist/${NAME}-Linux-aarch64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Linux-aarch64" >> "dist/${NAME}-sha256sum.txt"

.PHONY: build-dist-linux-armv7
build-dist-linux-armv7:
	@echo "Building ${COL_GREEN_BRIGHT}Linux-armv7l${COL_RESET}..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build ${FLAGS} -o "dist/${NAME}-Linux-armv7l" .
	@md5 -r "dist/${NAME}-Linux-armv7l" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Linux-armv7l" >> "dist/${NAME}-sha256sum.txt"

.PHONY: build-dist-darwin
build-dist-darwin: build-dist-darwin-amd64 build-dist-darwin-arm64

.PHONY: build-dist-darwin-amd64
build-dist-darwin-amd64:
	@echo "Building ${COL_GREEN_BRIGHT}Darwin-x86_64${COL_RESET}..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${FLAGS} -o "dist/${NAME}-Darwin-x86_64" .
	@md5 -r "dist/${NAME}-Darwin-x86_64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Darwin-x86_64" >> "dist/${NAME}-sha256sum.txt"

.PHONY: build-dist-darwin-arm64
build-dist-darwin-arm64:
	@echo "Building ${COL_GREEN_BRIGHT}Darwin-arm64${COL_RESET}..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build ${FLAGS} -o "dist/${NAME}-Darwin-arm64" .
	@md5 -r "dist/${NAME}-Darwin-arm64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Darwin-arm64" >> "dist/${NAME}-sha256sum.txt"

.PHONY: clean
clean:
	@echo "${COL_YELLOW_BRIGHT}Cleaning dist${COL_RESET}..."
	@rm -rf dist
