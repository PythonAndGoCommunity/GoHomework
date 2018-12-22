CURRENT_DIR = $(shell pwd)

check:
	go vet ./...

	goimports ./

	golint ./...

clean-binaries:
	rm server/server && rm client/client

build-server:
	go build -o server/server $(CURRENT_DIR)/server/server.go

build-client:
	go build -o client/client $(CURRENT_DIR)/client/client.go

clean:
	sudo docker system prune

build: clean
	sudo docker build -t "nonreldb" .

run: clean
	sudo docker run -d --net=host nonreldb

test: 
	echo "Running unit & integration tests"

	go test ./... -coverprofile coverage.out

	go tool cover -html=coverage.out




