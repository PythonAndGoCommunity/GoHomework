FROM golang:1.11

COPY . /go/src/gohomework

ENTRYPOINT ["/go/src/gohomework/server/server"]

EXPOSE 9090
