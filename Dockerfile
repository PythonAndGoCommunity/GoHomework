FROM golang:1.11-alpine3.8

EXPOSE 9090

RUN mkdir /NonRelDB/

ADD . /go/src/NonRelDB/

WORKDIR /go/src/NonRelDB/server

RUN go build *.go 

WORKDIR /go/src/NonRelDB/client

RUN go build *.go 

ENTRYPOINT [ "/go/src/NonRelDB/server/server" ]

