FROM golang:1.11

EXPOSE 9090

WORKDIR /usr/src/client

COPY ./ /usr/src/client

COPY ./docker-entrypoint.sh /usr/local/bin

RUN go get -u golang.org/x/lint/golint

RUN go get -u golang.org/x/tools/cmd/goimports

ENTRYPOINT docker-entrypoint.sh


