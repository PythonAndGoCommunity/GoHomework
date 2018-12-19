# Goredis Makefile

PROJECT_DIR_NAME = GoHomework

CLIENT_RUNFILE_PATH = client/runclient.go
SERVER_RUNFILE_PATH = server/runserver.go
CLIENT_TESTFILE_PATH = client/test/client_test.go
SERVER_TESTFILE_PATH = server/test/server_test.go

make_container:
	sudo docker build -t goredis-app .
	sudo docker run -d --rm -i --name goredis-app-running goredis-app

clear_dangling:
	sudo docker rmi $$(sudo docker images -f "dangling=true" -q)

enter_container:
	sudo docker exec -it goredis-app-running /bin/bash

stop_container:
	sudo docker stop goredis-app-running


check:
	go vet $(PROJECT_DIR_NAME)/server
	go vet $(PROJECT_DIR_NAME)/client

build:
	go build -o client/runclient $(CLIENT_RUNFILE_PATH)
	go build -o server/runserver $(SERVER_RUNFILE_PATH)

test:
	go test $(SERVER_TESTFILE_PATH)
	go test $(CLIENT_TESTFILE_PATH)

.PHONY: runserver runclient
runserver:
	./server/runserver

runclient:
	./client/runclient
