SERVER_PATH = ./src/server
SERVER_BIN = server

CLIENT_PATH = ./src/client
CLIENT_BIN = client

SEARCH_GOFILES =  find -not -path '*/vendor/*' -type f -name "*.go"

all: check client server docker run

clean:
	rm -rf $(SERVER_PATH)/server
	rm -rf $(CLIENT_PATH)/client

.PHONY: check
check:
	$(SEARCH_GOFILES) -exec go vet -s -l {} \;

client:
	go build -o $(CLIENT_PATH)/$(CLIENT_BIN) $(CLIENT_PATH)/client.go

server:
	go build -o $(SERVER_PATH)/$(SERVER_BIN) $(SERVER_PATH)/server.go	

.PHONY: docker
docker: server client 
	sudo docker build -t gohomework .

.PHONY: run
run:
	sudo docker run gohomework