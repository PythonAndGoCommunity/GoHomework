REPOSITORY_PATH := $(HOME)/Dev/Go/src/NonRelDB

VERSION := $(shell cat VERSION)

.PHONY: fmt
check:
	go vet NonRelDB/...

	goimports **/*.go

	golint **/*.go

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

	go test NonRelDB/... -coverprofile coverage.out

	go tool cover -html=coverage.out





