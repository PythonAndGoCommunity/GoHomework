SERVER_PATH := ./TinyRedisServer
CLIENT_PATH := ./TinyRedisClient
SERVER_FILE := ./TinyRedisServer/TinyRedisServer.go
CLIENT_FILE := ./TinyRedisClient/TinyRedisClient.go
SERVER_BIN := ./TinyRedisServer/TinyRedisServer
CLIENT_BIN := ./TinyRedisClient/TinyRedisClient

BECOME        := sudo -E
VERSION       := $(shell cat VERSION)
DOCKER_IMAGE   := tinyredis:$(VERSION)
DOCKER_BUILDER := golang:1.11
RUNNER =  docker run --rm -v $(CURDIR):/go/src/TinyRedis/$(SERVER_PATH)
RUNNER += $(DOCKER_ENVS) -w /go/src/TinyRedis/$(SERVER_PATH)

SEARCH_GOFILES =  find -not -path '*/vendor/*' -type f -name "*.go" ! -name "*_test*"
BUILDER = $(RUNNER) $(DOCKER_BUILDER)
PORT := 9090

.PHONY: tinyredis
tinyredis:
	docker build -t tinyredis .

.PHONY: build
build:
	go build -o $(SERVER_BIN) $(SERVER_FILE)
	go build -o $(CLIENT_BIN) $(CLIENT_FILE)
	
.PHONY: clean
clean:
	$(BECOME) $(RM) $(SERVER_PATH)/TinyRedisServer
	$(BECOME) $(RM) $(CLIENT_PATH)/TinyRedisClient

.PHONY: cleancontainer
cleancontainer:
	docker rm tinyredis

.PHONY: check
check:
	$(BECOME) $(BUILDER) sh -xc '\
		test -z "`$(SEARCH_GOFILES) -exec gofmt -s -l {} \;`" \
		&& test -z "`$(SEARCH_GOFILES) -exec go vet {} \;`"'

.PHONY: test
test:
	go test $(SERVER_PATH) -coverprofile coverage.out	
	go test $(CLIENT_PATH) -coverprofile coverage_client.out
	tail -n +2 coverage_client.out >> coverage.out
	rm coverage_client.out
	go tool cover -html coverage.out
