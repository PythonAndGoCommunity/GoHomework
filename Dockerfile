FROM golang:1.11

COPY src/server /usr/local/bin/server
COPY src/client /usr/local/bin/client

ENTRYPOINT ["/usr/local/bin/server"]

EXPOSE 9090
