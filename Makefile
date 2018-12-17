.DEFAULT_GOAL := build

DEV := dev
PROD := prod

APP_BASE_NAME := go-kvdb
APP_CLIENT_NAME := client
APP_SERVER_NAME := server

REPOSITORY_PATH := github.com/SiarheiKresik

DEV_GOPATH := /go
DEV_GOPATH_BIN := $(DEV_GOPATH)/bin
DEV_GOPATH_SRC := $(DEV_GOPATH)/src
DEV_WORKDIR := $(DEV_GOPATH_SRC)/$(REPOSITORY_PATH)/$(APP_BASE_NAME)


SRC_MOUNT := "$(PWD):$(DEV_GOPATH_SRC)/$(REPOSITORY_PATH)/$(APP_BASE_NAME)"
BIN_MOUNT := "$(PWD)/bin:$(DEV_GOPATH_BIN)"

DEV_IMAGE := $(APP_BASE_NAME)-${DEV}
PROD_IMAGE := $(APP_BASE_NAME)

BUILDER := docker build
RUNNER := docker run --rm
CHECKER := $(RUNNER) -v $(SRC_MOUNT) -w $(DEV_WORKDIR) $(DEV_IMAGE)

SEARCH_GOFILES = find -not -path '*/vendor/*' -type f -name "*.go"

DELIMITER="----------------------"
define print_target_name
	@echo $(DELIMITER)
	@echo $(1)
	@echo $(DELIMITER)
endef

### targets ###

build: build-dev-image build-dev build-prod-image

build-dev-image:
	$(call print_target_name, "Building an image with go tools for development...")
	$(BUILDER) -t $(DEV_IMAGE) --target $(DEV) .

build-dev:
	$(call print_target_name, "Compile binaries...")
	$(RUNNER) -v $(BIN_MOUNT) -v $(SRC_MOUNT) -w $(DEV_WORKDIR) $(DEV_IMAGE) sh -c "go install ./..."

build-prod-image:
	$(call print_target_name, "Building an image with server and client binaries")
	$(BUILDER) -t $(PROD_IMAGE) --target $(PROD) .

test:
	$(call print_target_name, "Run tests...")
	@echo "test are not implemented yet"

check: build-dev-image check-govet check-goimports check-golint

check-govet:
	$(call print_target_name, "Checks (go vet)...")
	$(CHECKER) sh -c "go tool vet -v ."

check-goimports:
	$(call print_target_name, "Checks (goimports)...")
	$(CHECKER) sh -c "$(SEARCH_GOFILES) -exec goimports {} \;"

check-golint:
	$(call print_target_name, "Checks (golint)...")
	$(CHECKER) sh -c "$(SEARCH_GOFILES) -exec golint {} \;"

run: run-server

run-server:
	$(call print_target_name, "Run server...")
	$(RUNNER) -p 9090:9090 $(APP_BASE_NAME):latest $(ARGS)

run-client:
	$(call print_target_name, "Run client...")
	$(RUNNER) -it --entrypoint /app/client $(APP_BASE_NAME):latest $(ARGS)

run-dev:
	$(RUNNER) -it -v $(SRC_MOUNT) -w $(DEV_WORKDIR) golang:1.11.1-alpine sh

prune:
	docker image prune
