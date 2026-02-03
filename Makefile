
SHELL = /bin/sh
PREFIX ?= /usr
NAME = frei
OUT_DIR = build

VERSION = $(shell git describe --tags --abbrev=0)
COMMIT = $(shell git rev-list -1 HEAD)

build: linux-amd64 linux-386 linux-arm darwin-amd64 darwin-arm64

linux-amd64:
	@echo "building linux-amd64"
	GOOS=linux GOARCH=amd64 \
		go build \
		-ldflags "-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT)" \
		-o "$(OUT_DIR)/$(NAME)-linux-amd64" .

linux-386:
	@echo "building linux-386"
	GOOS=linux GOARCH=386 \
		go build \
		-ldflags "-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT)" \
		-o "$(OUT_DIR)/$(NAME)-linux-386" .

linux-arm:
	@echo "building linux-arm"
	GOOS=linux GOARCH=arm \
		go build \
		-ldflags "-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT)" \
		-o "$(OUT_DIR)/$(NAME)-linux-arm" .

darwin-amd64:
	@echo "building darwin-amd64"
	GOOS=darwin GOARCH=amd64 \
		go build \
		-ldflags "-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT)" \
		-o "$(OUT_DIR)/$(NAME)-darwin-amd64" .

darwin-arm64:
	@echo "building darwin-arm64"
	GOOS=darwin GOARCH=arm64 \
		go build \
		-ldflags "-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT)" \
		-o "$(OUT_DIR)/$(NAME)-darwin-arm64" .

clean:
	$(RM) -r $(OUT_DIR)

.PHONY: build

