# Goredis Makefile

PROJECT_DIR_NAME = GoHomework

CLIENT_RUNFILE_PATH = client/runclient.go
SERVER_RUNFILE_PATH = server/runserver.go
SERVER_TESTFILE_PATH = server/test/server_test.go

check:
	go vet $(PROJECT_DIR_NAME)/server
	go vet $(PROJECT_DIR_NAME)/client

build:
	go build -o client/runclient $(CLIENT_RUNFILE_PATH)
	go build -o server/runserver $(SERVER_RUNFILE_PATH)

test:
	go test $(SERVER_TESTFILE_PATH)

.PHONY: runserver runclient
runserver:
	./server/runserver

runclient:
	./client/runclient
