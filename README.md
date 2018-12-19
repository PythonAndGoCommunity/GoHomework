# TinyRedis

*TinyRedis* is a very lightweight app, divided into two parts, *TinyRedisServer* & *TinyRedisClient*, which work together just like [Redis](https://redis.io), but only with 4 commands: `SET`, `GET`, `DEL` and `EXIT`. 

## Instruction
In order to launch the program, it's ought to do the following steps:
Get the app via following commands:
``` sh
git clone https://github.com/corvustristis/GoHomework
cd GoHomework
```

### Run via terminal
To do it, you have to install [Go](https://golang.org/) compiler. Then proceed to the following steps:
1) Launch *TinyRedisServer* from the folder of the same name:
```sh
./TinyRedisServer
```
There are options `-p` or `--port` for the choice of port and options `-m` or `--mode` for the choice of port.

2) Launch *TinyRedisClient* from another folder:
```sh
./TinyRedisClient
```
If you had chosen a port of your preferences in a previous step, don't forget to use it here as well. Optionally you can launch multiple clients at the same time, though you would be doing it at your own risk.

### Run via Docker
To proceed, you must have [Docker](https://www.docker.com/) installed. After that proceed to the following steps:

1) The project is already built and checked, but just in case you can delete binaries, compile them again, or check the code with `gofmt` and `go vet`, run some of the the following lines:
```sh
make clean build check
```

2) Build the project with:
```sh
make tinyredis
```

3) Launch server:
```sh
docker run -p 9090:9090 --name tinyredis tinyredis
```
There are options `-p` or `--port` for the choice of port and options `-m` or `--mode` for the choice of port. For example, if you have chosen custom port 3333, launch server via following:
```sh
docker run -p 3333:9090 --name tinyredis tinyredis --port=3333
```
If you have already lauched the container at least once, and there is a name conflict, do this: 
```sh
make cleancontainer
```

4) Launch bash for *TinyRedisClient*:
```sh
docker exec -it tinyredis bash
```

5) Finally, lauch the client itself:
```sh
./src/TinyRedis/TinyRedisClient/TinyRedisClient
```
Don't forget about your custom port here!

## Testing
Thought test files are stored in each folder separately, you can test the code via following instruction right from the main folder:
```sh
make test
```
After it, the browser with test data will open, and you would be able to choose either server, or client coverage.
