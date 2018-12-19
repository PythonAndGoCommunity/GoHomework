build:
	CGO_ENABLED=0 go install ./...
	docker build -t redislight .
test:
	go test -coverprofile=./coverage.profile ./...
	go tool cover --func=./coverage.profile > coverage.out
check:
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports

	${GOPATH}/bin/goimports -w .
	go vet ./...
	${GOPATH}/bin/golint ./...
run:
	CGO_ENABLED=0 go run ./cmd/server
