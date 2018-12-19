# Goredis

Goredis is a simple implementation of Redis on Golang.

## Getting started

### Requirements

All you need is a Docker app installed on your local computer.

### Installation

Clone or download this repository.

```git clone https://github.com/SamperMan44/GoHomework.git```

Build a docker image and run a Docker container by your own or using Makefile.

```make -f Makefile make_container```

Enter the terminal in your Docker container.

```make -f Makefile enter_container```

### Launching the server

Use Makefile to build and launch server with default params.

```
make -f Makefile build
make -f Makefile runserver
```

To get help on specified params, run server with argument `--help`. 

```
~/go$ cd src/GoHomework/server
~/go/src/GoHomework/server$ go run runserver.go --help
```

### Launching the client

The client is launched by analogy with the server. Simply use `runclient` instead of `runserver`.

### Checking and testing

Run Makefile commands `check` and `test`.

## Available features

To get the list of available commands, run client and type `HELP`.

```
localhost:9090> HELP
 Goredis is a simple implementation of Redis on Golang.
 Available commands: HELP SET GET DEL KEYS (UN)SUBSCRIBE PUBLISH EXIT STOP.
 To get more help on them, type any command with no arguments.
localhost:9090> 
```

Launch `runserver.go` or `runclient.go` with `--help` argument to get usage manual.

Server modes: `disk` (saves your changes), `memory` (keeps loaded database only in RAM).