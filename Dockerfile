# Goredis Dockerfile

FROM golang:latest

MAINTAINER Alexander Gutyra <gutyra13@gmail.com>

RUN mkdir -p /go/src/GoHomework

COPY . /go/src/GoHomework

EXPOSE 9090
