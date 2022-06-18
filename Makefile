
SHELL = /bin/sh
PREFIX ?= /usr
NAME = frei
OUT_DIR = build

TARGET_OS = linux
VERSION = $(shell git describe --tags --abbrev=0)
COMMIT = $(shell git rev-list -1 HEAD)

build: amd64 386 arm

amd64:
	@echo "building amd64"
	GOOS=$(TARGET_OS) GOARCH=amd64 \
		go build \
		-ldflags "-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT)" \
		-o "$(OUT_DIR)/$(NAME)-$(TARGET_OS)-amd64" .

386:
	@echo "building 386"
	GOOS=$(TARGET_OS) GOARCH=386 \
		go build \
		-ldflags "-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT)" \
		-o "$(OUT_DIR)/$(NAME)-$(TARGET_OS)-386" .

arm:
	@echo "building arm"
	GOOS=$(TARGET_OS) GOARCH=arm \
		go build \
		-ldflags "-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT)" \
		-o "$(OUT_DIR)/$(NAME)-$(TARGET_OS)-arm" .

clean:
	$(RM) -r $(OUT_DIR)

.PHONY: build

