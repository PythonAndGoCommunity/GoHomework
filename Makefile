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
DEV_WORKDIR := $(DEV_GOPATH_SRC)$(REPOSITORY_PATH)/$(APP_BASE_NAME)


SRC_MOUNT := "$(PWD):$(DEV_GOPATH_SRC)/$(REPOSITORY_PATH)/$(APP_BASE_NAME)"
BIN_MOUNT := "$(PWD)/bin:$(DEV_GOPATH_BIN)"

DEV_IMAGE := $(APP_BASE_NAME)-${DEV}
PROD_IMAGE := $(APP_BASE_NAME)

BUILDER := docker run --rm 
BUILDER += -v $(SRC_MOUNT) -v $(BIN_MOUNT)
BUILDER += -w $(DEV_WORKDIR) $(DEV_IMAGE)

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

build: build-dev-image build-dev build-prod-image clean

build-dev-image:
	$(call print_target_name, "Building an image with go tools for development...")
	$(IMAGE_BUILDER) -t $(DEV_IMAGE) --target $(DEV) .

build-dev:
	$(call print_target_name, "Compile binaries...")
	$(BUILDER) sh -c "go install -v ./..."

build-prod-image:
	$(call print_target_name, "Building an image with server and client binaries")
	$(IMAGE_BUILDER) -t $(PROD_IMAGE) --target $(PROD) .

test:
	$(call print_target_name, "Run tests...")
	@echo "test are not implemented yet"

check:
	$(BUILDER) sh -xc '\
		test -z "`$(SEARCH_GOFILES) -exec gofmt -s -l {} \;`" \
		&& test -z "`$(SEARCH_GOFILES) -exec golint {} \;`"'

che1ck:
	@echo "check"
	docker run \
	--rm \
	-v $(SRC_MOUNT) \
	$(DEV_IMAGE) \
	/bin/sh -c \
	"go vet src/*.go && goimports -w src/*.go && golint src/*.go"

run: run-server

run-server:
	$(call print_target_name, "Run server...")
	$(RUNNER) -p 9090:9090 $(APP_BASE_NAME):latest $(ARGS)

run-client:
	$(call print_target_name, "Run client...")
	$(RUNNER) --entrypoint /bin/client $(APP_BASE_NAME):latest $(ARGS)

clean:
	$(call print_target_name, "Cleaning up...")

prune:
	docker image prune