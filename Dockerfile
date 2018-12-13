# Goredis Dockerfile

FROM golang:latest

MAINTAINER Alexander Gutyra <gutyra13@gmail.com>

RUN mkdir -p /go/src/goredis

COPY . /go/src/goredis

EXPOSE 9090
