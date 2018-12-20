FROM golang:1.11

COPY . /Go/src/redis

WORKDIR /Go/src/
ENTRYPOINT ["/Go/pkg/redis"]
EXPOSE 9090
