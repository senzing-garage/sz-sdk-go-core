# Makefile for g2-sdk-go-base.

# Detect the operating system and architecture
include Makefile.osdetect

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

# "Simple expanded" variables (':=')

# PROGRAM_NAME is the name of the GIT repository.
PROGRAM_NAME := $(shell basename `git rev-parse --show-toplevel`)
MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIRECTORY := $(shell dirname $(MAKEFILE_PATH))
TARGET_DIRECTORY := $(MAKEFILE_DIRECTORY)/target
DOCKER_CONTAINER_NAME := $(PROGRAM_NAME)
DOCKER_IMAGE_NAME := senzing/$(PROGRAM_NAME)
DOCKER_BUILD_IMAGE_NAME := $(DOCKER_IMAGE_NAME)-build
BUILD_VERSION := $(shell git describe --always --tags --abbrev=0 --dirty  | sed 's/v//')
BUILD_TAG := $(shell git describe --always --tags --abbrev=0  | sed 's/v//')
BUILD_ITERATION := $(shell git log $(BUILD_TAG)..HEAD --oneline | wc -l | sed 's/^ *//')
GIT_REMOTE_URL := $(shell git config --get remote.origin.url)
GO_PACKAGE_NAME := $(shell echo $(GIT_REMOTE_URL) | sed -e 's|^git@github.com:|github.com/|' -e 's|\.git$$||' -e 's|Senzing|senzing|')
BIN_DIRECTORY := $(MAKEFILE_DIRECTORY)/bin
PATH := $(BIN_DIRECTORY):$(PATH)
GO_OSARCH = $(subst /, ,$@)
GO_OS = $(word 1, $(GO_OSARCH))
GO_ARCH = $(word 2, $(GO_OSARCH))

# set SQLite database variables
SENZING_TOOLS_DATABASE_PATH=$(TARGET_DIRECTORY)/sqlite/G2C.db
SENZING_TOOLS_DATABASE_URL ?= sqlite3://na:na@$(SENZING_TOOLS_DATABASE_PATH)

# Recursive assignment ('=')
CC = gcc

# Export environment variables.
.EXPORT_ALL_VARIABLES:

# -----------------------------------------------------------------------------
# Optionally include platform-specific settings and targets.
#  - Note: This is last because the "last one wins" when over-writing targets.
# -----------------------------------------------------------------------------
-include Makefile.$(OSTYPE)
-include Makefile.$(OSTYPE)_$(OSARCH)

# -----------------------------------------------------------------------------
# The first "make" target runs as default.
# -----------------------------------------------------------------------------

.PHONY: default
default: help

# -----------------------------------------------------------------------------
# Build
#  - The "build" target is implemented in Makefile.OS.ARCH files.
# -----------------------------------------------------------------------------

.PHONY: dependencies
dependencies:
	@go get -u ./...
	@go get -t -u ./...
	@go mod tidy


PLATFORMS := darwin/amd64 linux/amd64 windows/amd64
$(PLATFORMS):
	@echo Building $(TARGET_DIRECTORY)/$(GO_OS)-$(GO_ARCH)/$(PROGRAM_NAME)
	@mkdir -p $(TARGET_DIRECTORY)/$(GO_OS)-$(GO_ARCH) || true
	@GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build -o $(TARGET_DIRECTORY)/$(GO_OS)-$(GO_ARCH)/$(PROGRAM_NAME)


.PHONY: build-all $(PLATFORMS)
build-all: $(PLATFORMS)
	@mv $(TARGET_DIRECTORY)/windows-amd64/$(PROGRAM_NAME) $(TARGET_DIRECTORY)/windows-amd64/$(PROGRAM_NAME).exe

# -----------------------------------------------------------------------------
# Test
#  - The "test" target is implemented in Makefile.OS.ARCH files.
# -----------------------------------------------------------------------------

# -----------------------------------------------------------------------------
# Run
# -----------------------------------------------------------------------------

.PHONY: run
run:
	@go run main.go

# -----------------------------------------------------------------------------
# Utility targets
# -----------------------------------------------------------------------------

.PHONY: update-pkg-cache
update-pkg-cache:
	@GOPROXY=https://proxy.golang.org GO111MODULE=on \
		go get $(GO_PACKAGE_NAME)@$(BUILD_TAG)

.PHONY: setup
setup: 
	@mkdir -p $(shell dirname $(SENZING_TOOLS_DATABASE_PATH))
	@if [ ! -f $(SENZING_TOOLS_DATABASE_PATH) ]; then cp testdata/sqlite/G2C.db $(SENZING_TOOLS_DATABASE_PATH); fi
	
.PHONY: clean
clean:
	@go clean -cache
	@go clean -testcache
	@docker rm --force $(DOCKER_CONTAINER_NAME) 2> /dev/null || true
	@docker rmi --force $(DOCKER_IMAGE_NAME) $(DOCKER_BUILD_IMAGE_NAME) 2> /dev/null || true
	@rm -rf $(TARGET_DIRECTORY) || true
	@rm -f $(GOPATH)/bin/$(PROGRAM_NAME) || true
	@rm -rf $(shell dirname $(SENZING_TOOLS_DATABASE_PATH))
	@rm -rf /tmp/$(PROGRAM_NAME)

.PHONY: print-make-variables
print-make-variables:
	@$(foreach V,$(sort $(.VARIABLES)), \
		$(if $(filter-out environment% default automatic, \
		$(origin $V)),$(warning $V=$($V) ($(value $V)))))

# -----------------------------------------------------------------------------
# Help
# -----------------------------------------------------------------------------
.PHONY: help
help:
	@echo "Build $(PROGRAM_NAME) version $(BUILD_VERSION)-$(BUILD_ITERATION)".
	@echo "Makefile targets:"
	@$(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs
