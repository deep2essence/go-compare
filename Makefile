# detect operating system
ifeq ($(OS),Windows_NT)
    CURRENT_OS := Windows
else
    CURRENT_OS := $(shell uname -s)
endif

export GO111MODULE=on
export LOG_LEVEL=info
export CGO_ENABLED=1

#GOBIN
GOBIN = $(shell pwd)/build/bin
GO ?= latest

PACKAGES = $(shell go list ./... | grep -Ev 'vendor|importer')
COMMIT_HASH := $(shell git rev-parse --short HEAD)
GIT_BRANCH :=$(shell git branch 2>/dev/null | grep "^\*" | sed -e "s/^\*\ //")

build_tags = netgo
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

BUILD_FLAGS = -tags "$(build_tags)" -ldflags '-X github.com/deep2essence/gobin/version.GitCommit=${COMMIT_HASH} -X main.GitCommit=${COMMIT_HASH} -X main.DEBUGAPI=${DEBUGAPI} -X main.GitBranch=${GIT_BRANCH}'
BUILD_FLAGS_STATIC_LINK = -tags "$(build_tags)" -ldflags '-X github.com/deep2essence/gobin/version.GitCommit=${COMMIT_HASH} -X main.GitCommit=${COMMIT_HASH} -X main.DEBUGAPI=${DEBUGAPI} -X main.GitBranch=${GIT_BRANCH} -linkmode external -w -extldflags "-static"'

all: verify build

verify:
	@echo "verify modules"
	@go mod verify

update:
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -update

buildquick: go.sum
ifeq ($(CURRENT_OS),Windows)
	@echo BUILD_FLAGS=$(BUILD_FLAGS)
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bin/gobin.exe .
else
	@echo BUILD_FLAGS=$(BUILD_FLAGS)
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bin/gobin .
endif

build: unittest buildquick

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) .

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

# test part
test:
	@go test -v --vet=off $(PACKAGES)
	@echo $(PACKAGES)

unittest:
	@go test -v ./...

clear:
	@rm -rf *.lst

.PHONY: build install \
		test clean 
