FROM golang:1.11-alpine3.8

EXPOSE 9090

# Preparation stage.
RUN apk update && apk upgrade && apk add git

RUN go get golang.org/x/tools/cmd/goimports

RUN go get -u golang.org/x/lint/golint

RUN mkdir /NonRelDB/

ADD . /go/src/NonRelDB/


# Build stage.
WORKDIR /go/src/NonRelDB/server

RUN go build *.go 

WORKDIR /go/src/NonRelDB/client

RUN go build *.go 

# Check stage.
WORKDIR /go/src/NonRelDB

RUN go vet **/*.go

RUN goimports **/*.go

RUN golint **/*.go

# Entrypoint bind.
ENTRYPOINT [ "/go/src/NonRelDB/server/server" ]

