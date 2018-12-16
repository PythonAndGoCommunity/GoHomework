SERVER_PATH = ./src/server
SERVER_BIN = server

CLIENT_PATH = ./src/client
CLIENT_BIN = client

SEARCH_GOFILES := $(shell find -type f -name "*.go")

build: clean client server docker

.PHONY: test
test:
	go test -coverprofile=coverage.out $(SERVER_PATH) $(CLIENT_PATH)

.PHONY: check
check: 
	goimports -e -l $(SEARCH_GOFILES)
	golint -set_exit_status $(SEARCH_GOFILES)
	go vet $(SERVER_PATH)
	go vet $(CLIENT_PATH)

.PHONY: run
run: buid
	sudo docker run -d gohomework
	#docker exec -i $(docker ps -qf ancestor=gohomework --last=1) /usr/local/bin/client

clean:
	rm -rf $(SERVER_PATH)/$(SERVER_BIN)
	rm -rf $(CLIENT_PATH)/$(CLIENT_BIN)

client:
	go build -o $(CLIENT_PATH)/$(CLIENT_BIN) $(CLIENT_PATH)/client.go

server:
	go build -o $(SERVER_PATH)/$(SERVER_BIN) $(SERVER_PATH)/server.go	

.PHONY: docker
docker: server client 
	sudo docker build -t "gohomework" .