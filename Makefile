SRC_PATH = ./src
SERVER_PATH = ./src/server
SERVER_BIN = server
CLIENT_PATH = ./src/client
CLIENT_BIN = client

DOCKER_BUILDER := golang:1.11
BECOME        := sudo -E
VERSION       := $(shell cat VERSION)
DOCKER_IMAGE   := gohomework:$(VERSION)
RUNNER =  docker run --rm -v $(CURDIR):/go/src/gohomework/$(SERVER_PATH)
RUNNER += $(DOCKER_ENVS) -w /go/src/gohomework/$(SERVER_PATH)

BUILDER = $(RUNNER) $(DOCKER_BUILDER)
PORT := 9090

.PHONY: test
test:
	go test  -coverprofile coverage.out -v ./...

.PHONY: check
check: 
	goimports -e -l $(SRC_PATH)
	golint -set_exit_status $(SRC_PATH)
	go vet $(SERVER_PATH)
	go vet $(CLIENT_PATH)


.PHONY: gohomework
tinyredis:
	docker build -t gohomework .

.PHONY: build
build:
	go build -o $(SERVER_BIN) $(SERVER_FILE)
	go build -o $(CLIENT_BIN) $(CLIENT_FILE)
	
.PHONY: clean
clean:
	$(BECOME) $(RM) $(SERVER_PATH)/server
	$(BECOME) $(RM) $(CLIENT_PATH)/client

.PHONY: cleancontainer
cleancontainer:
	docker rm gohomework