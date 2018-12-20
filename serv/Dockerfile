FROM golang:1.11

EXPOSE 9090

WORKDIR /usr/src/serv

COPY ./ /usr/src/serv

COPY ./docker-entrypoint.sh /usr/local/bin

RUN go get -u golang.org/x/lint/golint

RUN go get -u golang.org/x/tools/cmd/goimports

ENTRYPOINT docker-entrypoint.sh


