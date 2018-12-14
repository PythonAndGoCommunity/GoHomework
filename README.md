# TinyRedis

*TinyRedis* is a very lightweight app, divided into two parts, *TinyRedisServer* & *TinyRedisClient*, which work together just like [Redis](https://redis.io), but only with 4 commands: `SET`, `GET`, `DEL` and `EXIT`. 

### Instruction
In order to launch the program, it's ought to do the following steps:
1) Get the app via following commands:
``` sh
git clone https://github.com/corvustristis/GoHomework
cd GoHomework
```
2) Launch *TinyRedisServer* from the folder of the same name:
```sh
./TinyRedisServer
```
There are options `-p` or `--port` for the choice of port and options `-m` or `--mode` for the choice of port.

3) Launch *TinyRedisClient* from another folder:
```sh
./TinyRedisClient
```
If you had chosen a port of your preferences in a previous step, don't forget to use it here as well. Optionally you can launch multiple clients at the same time, though you would be doing it at your own risk.

### Testing
Test data is avaliable in each folder separately. In order to test the code, run for each part of the application:
```sh
make test
```

### Docker
*Adequate docker files unfortunatelly are not yet avaliable.*
