VERSION := $(shell cat VERSION)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%S)
COMMIT := $(shell git log --pretty=format:'%h' -n 1)
BIN_PATH := bin/duclean

init_dir:
	mkdir -p bin/
	mkdir -p releases/

get_deps:
	git submodule init
	git submodule update

test:
	go test

build: init_dir
	go build -o $(BIN_PATH) \
		-ldflags "-X main.version=$(VERSION) \
			-X main.buildTime=$(BUILD_TIME) \
			-X main.commit=$(COMMIT)" \
		*.go

push_tag:
	git tag v$(VERSION)
	git push origin v$(VERSION)

archive: init_dir build
	gzip -c $(BIN_PATH) > releases/duclean-v$(VERSION).gz
