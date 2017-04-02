VERSION := $(shell cat VERSION)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%S)
COMMIT := $(shell git log --pretty=format:'%h' -n 1)

init_dir:
	mkdir -p bin/

build: init_dir
	go build -o bin/duclean \
		-ldflags "-X main.version=$(VERSION) \
			-X main.buildTime=$(BUILD_TIME) \
			-X main.commit=$(COMMIT)" \
		*.go
