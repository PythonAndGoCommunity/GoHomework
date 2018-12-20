SRC_PATH = ./src
SERVER_PATH = ./src/server
SERVER_BIN = server
CLIENT_PATH = ./src/client
CLIENT_BIN = client

.PHONY: test
test:
	go test  -coverprofile coverage.out -v ./...
.PHONY: check
check: 
	goimports -e -l $(SRC_PATH)
	golint -set_exit_status $(SRC_PATH)
	go vet $(SERVER_PATH)
	go vet $(CLIENT_PATH)
