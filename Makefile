REPOSITORY_PATH := $(HOME)/Dev/Go/src/NonRelDB

VERSION := $(shell cat VERSION)

.PHONY: fmt
check:
	go vet $(REPOSITORY_PATH)/server
	go vet $(REPOSITORY_PATH)/client
	go vet $(REPOSITORY_PATH)/log
	go vet $(REPOSITORY_PATH)/util/collection
	go vet $(REPOSITORY_PATH)/util/file
	go vet $(REPOSITORY_PATH)/util/json
	go vet $(REPOSITORY_PATH)/util/regex
	go vet $(REPOSITORY_PATH)/util/sync

	goimports $(REPOSITORY_PATH)/server
	goimports $(REPOSITORY_PATH)/client
	goimports $(REPOSITORY_PATH)/log
	goimports $(REPOSITORY_PATH)/util/collection
	goimports $(REPOSITORY_PATH)/util/file
	goimports $(REPOSITORY_PATH)/util/json
	goimports $(REPOSITORY_PATH)/util/regex
	goimports $(REPOSITORY_PATH)/util/sync


	golint $(REPOSITORY_PATH)/server
	golint $(REPOSITORY_PATH)/client
	golint $(REPOSITORY_PATH)/log
	golint $(REPOSITORY_PATH)/util/collection
	golint $(REPOSITORY_PATH)/util/file
	golint $(REPOSITORY_PATH)/util/json
	golint $(REPOSITORY_PATH)/util/regex
	golint $(REPOSITORY_PATH)/util/sync

clean:
	rm server/server && rm client/client

build-server:
	go build -o server/server $(REPOSITORY_PATH)/server/server.go

build-client:
	go build -o client/client $(REPOSITORY_PATH)/client/client.go

build-container:
	sudo docker build -t "nonreldb" .

run:
	sudo docker run --net=host nonreldb

test: 
	echo "Running unit & integration tests"
	go test $(REPOSITORY_PATH)/server/handler -coverprofile cover.out
	go test $(REPOSITORY_PATH)/server/storage/inmemory -coverprofile cover.out
	go test $(REPOSITORY_PATH)/util/collection -coverprofile cover.out
	go test $(REPOSITORY_PATH)/util/file -coverprofile cover.out
	go test $(REPOSITORY_PATH)/util/json -coverprofile cover.out
	go test $(REPOSITORY_PATH)/util/regex -coverprofile cover.out
	go test $(REPOSITORY_PATH)/util/sync -coverprofile cover.out




