
ARG GO_VERSION=1.11.1
FROM golang:${GO_VERSION}-alpine AS dev
RUN apk add git
RUN go get golang.org/x/lint/golint
RUN go get golang.org/x/tools/cmd/goimports

FROM alpine AS prod
COPY ./bin/client /bin/
COPY ./bin/server /bin/
EXPOSE 9090
ENTRYPOINT ["/bin/server"]