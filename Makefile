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

BUILDER := docker run --rm 
BUILDER += -v $(SRC_MOUNT)
BUILDER += -w $(DEV_WORKDIR)

#  -v $(BIN_MOUNT)
IMAGE_BUILDER := docker build
RUNNER := docker run --rm -it

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
	$(IMAGE_BUILDER) -t $(DEV_IMAGE) --target $(DEV) .

build-dev:
	$(call print_target_name, "Compile binaries...")
	$(BUILDER) -v $(BIN_MOUNT) $(DEV_IMAGE) sh -c "go install ./..."

build-prod-image:
	$(call print_target_name, "Building an image with server and client binaries")
	$(IMAGE_BUILDER) -t $(PROD_IMAGE) --target $(PROD) .

test:
	$(call print_target_name, "Run tests...")
	@echo "test are not implemented yet"

check: build-dev-image
	$(BUILDER) $(DEV_IMAGE) sh -xc "\
		go version && \
		$(SEARCH_GOFILES) -exec gofmt -s -l {} \; && \
		$(SEARCH_GOFILES) -exec golint {} \;"

run: run-server

run-server:
	$(call print_target_name, "Run server...")
	$(RUNNER) -p 9090:9090 $(APP_BASE_NAME):latest $(ARGS)

run-client:
	$(call print_target_name, "Run client...")
	$(RUNNER) --entrypoint /app/client $(APP_BASE_NAME):latest $(ARGS)

prune:
	docker image prune