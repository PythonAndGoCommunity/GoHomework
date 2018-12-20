# go-kvdb

Go implementation of the server-client solution for storing key-value data.

## Requirements

-   [make](https://www.gnu.org/software/make/)
-   [docker](https://www.docker.com/)  
    Note: in order to run docker commands in the makefile targets without sudo you have to follow the next [guide](https://docs.docker.com/install/linux/linux-postinstall/).

## Getting source files

If you have `go` tool installed:

-   get source files: `go get github.com/SiarheiKresik/go-kvdb`
-   `cd $GOPATH/src/github.com/SiarheiKresik/go-kvdb`

Or get source files by cloning repository directly:

-   get source files: `git clone https://github.com/SiarheiKresik/go-kvdb.git`
-   `cd go-kvdb`

## Building

Run `make`. This builds the `go-kvdb` docker container with the server and client binaries.

## Running

Run server in a docker container:

-   `make run [ARGS="<args>"]` or
-   `docker run --rm -p 9090:9090 go-kvdb:latest [ARG...]`

Run client in a docker container:

-   `make run-client [ARGS="<args>"]` or
-   `docker run --rm -it --entrypoint /app/client go-kvdb:latest [ARG...]`

## Other make targets

-   `make check`—runs code linting and verification
-   `make test`—runs tests for the project files
-   `make coverage`—measures test coverage, saves result to the `coverage.out` file

## Server

Starting options:

-   `-m`, `--mode`
    The possible storage option (_disk_—stores database to disk, _memory_—database runs only in memory)
-   `-p`, `--port`
    The port for listening on (default—_9090_)

## Client

Starting options:

-   `-p`, `--port`
    The port to connect to the server (default—_9090_)
-   `-h`, `--host`
    The host to connect to the server (default—_127.0.0.1_)

## Supported commands

-   **SET** - updates a key at a time with the given value
-   **GET** - returns the value of a key
-   **DEL** - removes a key
-   **KEYS** - returns all keys matching pattern; pattern could include '\*' symbol which matches zero or
    more characters

## References

1. https://systembash.com/a-simple-go-tcp-server-and-tcp-client/
2. https://www.oreilly.com/ideas/building-messaging-in-go-network-clients
3. https://medium.com/@kpbird/golang-serialize-struct-using-gob-part-1-e927a6547c00
