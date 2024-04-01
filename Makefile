#!/usr/bin/make -f

BRANCH         := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT         := $(shell git log -1 --format='%H')
BUILD_DIR      ?= $(CURDIR)/build
DIST_DIR       ?= $(CURDIR)/dist
LEDGER_ENABLED ?= true
TM_VERSION     := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')
DOCKER         := $(shell which docker)
PROJECT_NAME   := ssc
HTTPS_GIT      := https://github.com/sagaxyz/ssc.git

ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

###############################################################################
##                                   Build                                   ##
###############################################################################

build_tags = netgo

ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
	  LEDGER_ENABLED=false
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=ssc \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=sscd \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TM_VERSION)

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

build: go.sum
	@echo "Building..."
	go build -mod=readonly $(BUILD_FLAGS) -o $(BUILD_DIR)/ ./...

build-no_cgo:
	@echo "Building static binary with no CGO nor GLIBC dynamic linking..."
	CGO_ENABLED=0 CGO_LDFLAGS="-static" $(MAKE) build

# Linux Targets
build-linux-amd64:
	LEDGER_ENABLED=$(LEDGER_ENABLED) GOOS=linux GOARCH=amd64 $(MAKE) build

build-linux-386:
	LEDGER_ENABLED=$(LEDGER_ENABLED) GOOS=linux GOARCH=386 $(MAKE) build

build-linux-arm:
	LEDGER_ENABLED=$(LEDGER_ENABLED) GOOS=linux GOARCH=arm $(MAKE) build

build-linux-arm64:
	LEDGER_ENABLED=$(LEDGER_ENABLED) GOOS=linux GOARCH=arm64 $(MAKE) build

# macOS (Darwin) Targets
build-darwin-amd64:
	LEDGER_ENABLED=$(LEDGER_ENABLED) GOOS=darwin GOARCH=amd64 $(MAKE) build

build-darwin-arm64: # For M1 Macs
	LEDGER_ENABLED=$(LEDGER_ENABLED) GOOS=darwin GOARCH=arm64 $(MAKE) build

# Windows Targets
build-windows-amd64:
	LEDGER_ENABLED=$(LEDGER_ENABLED) GOOS=windows GOARCH=amd64 $(MAKE) build

build-windows-386:
	LEDGER_ENABLED=$(LEDGER_ENABLED) GOOS=windows GOARCH=386 $(MAKE) build

install: go.sum
	@echo "Installing..."
	LEDGER_ENABLED=$(LEDGER_ENABLED) go install -mod=readonly $(BUILD_FLAGS) ./...

go-mod-cache: go.sum
	@echo "Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "Ensure dependencies have not been modified"
	@go mod verify

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)/**  $(DIST_DIR)/**

.PHONY: install build build-linux-amd64 build-linux-386 build-linux-arm build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-windows-386 clean
