################ Preparation ####################
REGISTRY := nick1chak
VERSION := $(shell cat VERSION)
DOCKER_IMAGE_SERV := $(REGISTRY)/redis:$(VERSION)
DOCKER_IMAGE_CLI := $(REGISTRY)/red_client:$(VERSION)

REPOSITORY_PATH := /usr/src/redis
SERVER_PATH := serv
CLIENT_PATH := client

DOCKER_BUILDER := golang:1.11
DOCKER_RUNNER =  docker run --rm -v $(CURDIR):$(REPOSITORY_PATH)
DOCKER_RUNNER += -w $(REPOSITORY_PATH)/
################ End Preparation ####################



################ Binary Target ####################
.PHONY: build
build: build_cli build_serv
	docker build -t redis:build .
	docker run -it --name redis_build redis:build
	./serv/serv &
	./client/client

.PHONY: build_cli
build_cli:
	$(DOCKER_RUNNER)$(CLIENT_PATH) $(DOCKER_BUILDER) go build

.PHONY: build_serv
build_serv:
	$(DOCKER_RUNNER)$(SERVER_PATH) $(DOCKER_BUILDER) go build

############### Docker Target ####################
.PHONY: run
run: build_serv
	docker build -t redis:default ./serv/
	docker run --name redis_default redis:default

.PHONY: build_check
build_check:
	docker build -t redis:check .
################# Testing ####################
.PHONY: check
check: build_check
	docker run -it --name redis_check redis:check /bin/bash
	echo go vet: client
	go vet ./client/

.PHONY: test
test:
	$(DOCKER_RUNNER)$(SERVER_PATH) $(DOCKER_BUILDER) go test -cover -coverprofile=coverage.out
