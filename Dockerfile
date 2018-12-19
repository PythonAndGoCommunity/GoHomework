# https://hub.docker.com/_/golang/
FROM golang:1.11.3-alpine3.8

RUN mkdir -p /go/src/github.com/ITandElectronics/GoHomework
WORKDIR /go/src/github.com/ITandElectronics/GoHomework
COPY . .

RUN CGO_ENABLED=0 go test -v ./...

RUN CGO_ENABLED=0 go install -v ./...

FROM alpine:3.8

WORKDIR /root/
COPY --from=0 /go/bin/client .
COPY --from=0 /go/bin/server .

CMD ["./server"]
