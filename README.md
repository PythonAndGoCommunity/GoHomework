## Redislight 
### Simplified version of a redis. Supports only GET, SET and DEL commands

## Usage
### Pre Requisites
 - Docker  

Redis light consists of two application: client and server. They both bundled into a single docker image which can be built by running following command:
```bash
$ docker build -t redislight .
```
### Server
In order to run server:
```bash
$ docker run -p 9090:9090 redislight
```
Server supports following options:
```bash
./server --help
Usage of server:
  --mode, -m string
          Storage options. One of [disk] (default "disk")
  --port, -p int
          Port to listen on (default 9090)
```

### Client
In order to run client:
```bash
$ docker run redislight ./client
```

Client supports following options:
```bash 
./client --help
Usage of client:
  --host, -h string
          Remote server address (default "127.0.0.1")
  --port, -p int
          Remote server port (default 9090)
```

### Supported commands
**GET** *key* - return value associated with proveded *key* or 'key is not exists' error  
**SET** *key* *value*  - create or update *value* associated with the *key*  
**DEL** *key* -  remove value associated with the *key*. If *key* is not exists, it will return 'not exists error'


## Development
### Pre Requisites
 - Go  >= 1.11
 - Docker
 - Make

Clone this repository into GOPATH/src/{repository}/redislight:
```bash 
$ git clone .../redislight.git
```  
Run linters:
```bash
$ make check
```
Run tests(also will produce code coverage in coverage.out file):
```bash
$ make tests
```
Build:
```bash
$ make build
```

Run server:
```bash
$ make run
```