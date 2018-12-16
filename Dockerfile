FROM golang:latest

EXPOSE 9090

COPY src/server/server /usr/local/bin/server

COPY src/client/client /usr/local/bin/client

ENTRYPOINT ["/usr/local/bin/server"]
