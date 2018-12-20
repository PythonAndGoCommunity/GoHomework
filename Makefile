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

.PHONY: build_cli
build_cli:
	$(DOCKER_RUNNER)$(CLIENT_PATH) $(DOCKER_BUILDER) go build

.PHONY: build_serv
build_serv:
	$(DOCKER_RUNNER)$(SERVER_PATH) $(DOCKER_BUILDER) go build

############### Docker Target ####################
.PHONY: run
run: build
	docker run -it --name redis_build redis:build

.PHONY: build_check
build_check:
	docker build -t redis:checkSR ./$(SERVER_PATH)
	docker build -t redis:checkCL ./$(CLIENT_PATH)
################# Testing ####################
.PHONY: check
check: build_check
	docker run --rm --name redis_checkSR redis:checkSR
	docker run --rm --name redis_checkCL redis:checkCL

.PHONY: test
test:
	$(DOCKER_RUNNER)$(SERVER_PATH) $(DOCKER_BUILDER) go test -cover -coverprofile=coverage.out

.PHONY: clean
clean:
	docker rm -f  redis_build &>/dev/null
	docker rmi -f redis:build &>/dev/null
	docker rmi -f redis:checkSR &>/dev/null
	docker rmi -f redis:checkCL &>/dev/null
