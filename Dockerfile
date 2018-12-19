FROM golang:latest

COPY . /go/src/TinyRedis

ENTRYPOINT ["/go/src/TinyRedis/TinyRedisServer/TinyRedisServer"]

EXPOSE 9090
