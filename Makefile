# Copyright wi-cuckoo@github.com 2020

SHELL := /bin/bash

get-command = $(shell which="$$(which $(1) 2> /dev/null)" && if [[ ! -z "$$which" ]]; then printf %q "$$which"; fi)

GO 			:= $(call get-command,go)

VERSION  := $(shell git describe --exact-match --tags 2>/dev/null)
REVISION := $(shell git rev-parse --short HEAD)
GOFILES  ?= $(shell git ls-files '*.go' | grep -v vendor/)
GOFMT 	 ?= $(shell gofmt -l -s $(GOFILES))

ifndef VERSION
	VERSION = 0.1.2
endif

LDFLAGS := $(LDFLAGS) -X main.revision=$(REVISION) -X main.version=$(VERSION)

PACKAGE_NAME ?= $(shell ls release/ | grep $(VERSION))
# package env
PACKAGE_ROOT := /tmp/$(shell date +%s)_cc_pkg


all: clean test-all build package

.PHONY: fmt
fmt:
	@gofmt -s -w $(GOFILES)

.PHONY: fmtcheck
fmtcheck:
	@echo $(GOFMT)
	@if [ ! -z "$(GOFMT)" ]; then \
		echo "[ERROR] gofmt has found errors in the following files:"  ; \
		echo "$(GOFMT)" ; \
		echo "" ;\
		echo "Run make fmt to fix them." ; \
		exit 1 ;\
	fi

vet:
	@echo 'go vet $$(go list ./... | grep -v ./vendor/)'
	@go vet $$(go list ./... | grep -v ./vendor/) ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "go vet has found suspicious constructs. Please remediate any reported errors"; \
		echo "to fix them before submitting code for review."; \
		exit 1; \
	fi

test-all: fmt vet
	go test ./...

# build binary
build:
	@echo "Building totally static linux binary of server side ......"
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" -o ./bin/goaheadd ./cmd/server

build-cli:
	@echo "Building totally static linux binary of client side ......"
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" -o ./bin/goahead ./cmd/cli


# package .rpm / .deb
.PHONY: package
package: build build-cli
	# create tmp dir
	mkdir -p $(PACKAGE_ROOT)/usr/sbin/
	# copy file
	cp -p ./bin/goaheadd $(PACKAGE_ROOT)/usr/sbin
	cp -p ./bin/goahead $(PACKAGE_ROOT)/usr/sbin
	fpm -t rpm \
		-s dir \
		-C $(PACKAGE_ROOT) \
		-p ./release \
		--name "goahead" \
		--description "a tool like supervisor" \
		--url github.com/wi-cuckoo/goahead \
		--maintainer wi-cuckoo@github.com \
		--version $(VERSION) \
		--iteration stable \
		.

.PHONY: clean
clean:
	rm -rf ./bin/*
	rm -rf /release/*