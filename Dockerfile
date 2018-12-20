FROM golang:1.11

COPY ./src/server/server /go/src/server
COPY ./src/client/client /go/src/client
RUN go get -u golang.org/x/lint/golint
RUN go get golang.org/x/tools/cmd/goimports
RUN go get github.com/golang/go/src/cmd/vet

WORKDIR /go/src/
RUN go build /go/src/server.go
ENTRYPOINT ["/go/src/server"]
EXPOSE 9090
