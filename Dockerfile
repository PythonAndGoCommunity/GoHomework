FROM golang:1.11

COPY ./server/server_bin /go/src/
COPY ./server /go/src/server
COPY ./client /go/src/client
RUN go get -u golang.org/x/lint/golint && \
go get golang.org/x/tools/cmd/goimports && \
go get github.com/golang/go/src/cmd/vet

WORKDIR /go/src/
ENTRYPOINT ["/go/src/server_bin"]
EXPOSE 9090